package main

func (s *Store) updateFactory(id, rawName string) (Factory, error) {
	name, err := parseName(rawName)
	if err != nil {
		return Factory{}, err
	}
	var factory Factory
	err = s.db.QueryRow("UPDATE factories SET name = ? WHERE id = ? RETURNING id,name,workspace_root,repository_path", name, id).Scan(&factory.ID, &factory.Name, &factory.WorkspaceRoot, &factory.RepositoryPath)
	return factory, err
}
func (s *Store) deleteFactory(id string) error {
	result, err := s.db.Exec("DELETE FROM factories WHERE id = ? AND NOT EXISTS (SELECT 1 FROM user_stories WHERE factory_id = ?)", id, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil || affected > 0 {
		return err
	}
	var exists int
	if err = s.db.QueryRow("SELECT 1 FROM factories WHERE id = ?", id).Scan(&exists); err != nil {
		return err
	}
	return ErrFactoryNotEmpty
}
