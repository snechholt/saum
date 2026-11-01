package main

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestCodemodeDiscoveryAndExecution(t *testing.T) {
	store := testStore(t)
	session := testQueuedSession(t, store)
	if err := store.startSession(session.ID, t.TempDir(), "prompt", "secret"); err != nil {
		t.Fatal(err)
	}
	found := searchFunctions("saum.submit_item")
	if len(found) != 1 || found[0].Name != "saum.submit_item" {
		t.Fatalf("unexpected discovery: %+v", found)
	}
	code := `const item=saum.submit_item({key:"one",title:"Deliver sign-in",description:"Allow valid users to sign in.",acceptanceCriteria:["Valid users can sign in"]}); return {id:item.id,plan:saum.complete_plan()};`
	result, err := store.executeCode(context.Background(), session.ID, code)
	if err != nil || result == nil {
		t.Fatalf("execute failed: %+v %v", result, err)
	}
	items, err := store.items(session.ID)
	if err != nil || len(items) != 1 {
		t.Fatalf("missing item: %+v %v", items, err)
	}
	for _, code := range []string{`require("fs")`, `fetch("http://example.com")`, `process.exit()`, `saum.submit_item({key:"two",title:"Second",description:"Work",acceptanceCriteria:["Done"]})`} {
		if _, err := store.executeCode(context.Background(), session.ID, code); err == nil {
			t.Fatalf("accepted forbidden code: %s", code)
		}
	}
}
func TestCodemodeInterruptsRunawayCode(t *testing.T) {
	store := testStore(t)
	session := testQueuedSession(t, store)
	if err := store.startSession(session.ID, t.TempDir(), "prompt", "secret"); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	started := time.Now()
	_, err := store.executeCode(ctx, session.ID, "while(true) {}")
	if err == nil || time.Since(started) > time.Second {
		t.Fatalf("code not interrupted promptly: %v", err)
	}
	events, err := store.events(session.ID, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) < 2 || !strings.Contains(string(events[len(events)-1].Data), "error") {
		t.Fatal("execution failure was not saved")
	}
}
