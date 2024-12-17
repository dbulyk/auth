package user

import "context"

// DeleteUserServ является сервисной прослойкой для удаления пользователя
func (s *serv) DeleteUserServ(ctx context.Context, id int64) (err error) {
	err = s.authRepository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
