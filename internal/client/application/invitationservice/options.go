package invitationservice

type Option func(s *service) error

//func WithKeystoreRepository(repo invitation.Repository) Option {
//	return func(s *service) error {
//		s.invitations = repo
//		return nil
//	}
//}
