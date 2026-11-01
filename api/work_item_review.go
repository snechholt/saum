package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

func (s *Store) workItemReviews(itemID string) ([]ExecutiveReview, error) {
	rows, err := s.db.Query("SELECT id,work_item_id,goal_session_id,summary_session_id,summary,verification,risks,decision,feedback,created_at,reviewed_at FROM work_item_reviews WHERE work_item_id=? ORDER BY rowid", itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reviews := []ExecutiveReview{}
	for rows.Next() {
		var review ExecutiveReview
		var risks string
		if err := rows.Scan(&review.ID, &review.WorkItemID, &review.GoalSessionID, &review.SummarySessionID, &review.Summary, &review.Verification, &risks, &review.Decision, &review.Feedback, &review.CreatedAt, &review.ReviewedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(risks), &review.Risks); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, rows.Err()
}

type WorkItemDecision struct {
	ReviewID string `json:"reviewId"`
	Feedback string `json:"feedback"`
}

func (s *Store) decideWorkItem(itemID string, input WorkItemDecision, approve bool) (*Session, error) {
	feedback := ""
	var err error
	if !approve {
		feedback, err = progressSummary(input.Feedback)
		if err != nil {
			return nil, err
		}
	} else if compact(input.Feedback) != "" {
		return nil, fmt.Errorf("send feedback instead of approving with feedback")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var state, planID, latestReview, summarySessionID string
	var decision ReviewStatus
	var sessionState SessionStatus
	if err := tx.QueryRow(`SELECT work.status,work.plan_session_id,review.id,review.decision,review.summary_session_id,run.status FROM work_items AS work
 JOIN work_item_reviews AS review ON review.id=(SELECT id FROM work_item_reviews WHERE work_item_id=work.id ORDER BY rowid DESC LIMIT 1)
 JOIN sessions AS run ON run.id=review.summary_session_id WHERE work.id=?`, itemID).Scan(&state, &planID, &latestReview, &decision, &summarySessionID, &sessionState); err != nil {
		return nil, err
	}
	if state != "awaiting_review" || latestReview != input.ReviewID || decision != ReviewPending || sessionState != SessionCompleted {
		return nil, fmt.Errorf("executive summary is no longer available for review")
	}
	nextDecision := ReviewChangesRequested
	if approve {
		nextDecision = ReviewApproved
	}
	if _, err := tx.Exec("UPDATE work_item_reviews SET decision=?,feedback=?,reviewed_at=? WHERE id=?", nextDecision, feedback, time.Now().UTC().Format(time.RFC3339Nano), latestReview); err != nil {
		return nil, err
	}
	var next *Session
	if !approve {
		session, err := queueWorkItemSession(tx, itemID, ImplementationSession)
		if err != nil {
			return nil, err
		}
		next = &session
	} else {
		if _, err := tx.Exec("UPDATE work_items SET status='approved' WHERE id=?", itemID); err != nil {
			return nil, err
		}
		var nextID, nextState string
		err := tx.QueryRow("SELECT id,status FROM work_items WHERE plan_session_id=? AND status<>'approved' ORDER BY position LIMIT 1", planID).Scan(&nextID, &nextState)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == nil {
			if nextState != "pending" {
				return nil, fmt.Errorf("next work item is not pending")
			}
			session, err := queueWorkItemSession(tx, nextID, ImplementationSession)
			if err != nil {
				return nil, err
			}
			next = &session
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return next, nil
}
