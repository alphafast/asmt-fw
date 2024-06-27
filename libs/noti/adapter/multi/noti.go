package multi

import (
	"context"

	"github.com/alphafast/asmt-fw/libs/domain/noti"
	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
	"github.com/pkg/errors"
)

type MultiNotiAdapter struct {
	emailAdapter []noti.NotiAdapter
	pushAdapter  []noti.NotiAdapter
}

type withOption = func(m *MultiNotiAdapter)

func WithBackupEmailAdapter(emailAdapter noti.NotiAdapter) withOption {
	return func(m *MultiNotiAdapter) {
		m.emailAdapter = append(m.emailAdapter, emailAdapter)
	}
}

func WithBackupPushAdapter(emailAdapter noti.NotiAdapter) withOption {
	return func(m *MultiNotiAdapter) {
		m.emailAdapter = append(m.emailAdapter, emailAdapter)
	}
}

func NewNotiAdapter(emailAdapter noti.NotiAdapter, pushAdapter noti.NotiAdapter, opts ...withOption) *MultiNotiAdapter {
	adapter := &MultiNotiAdapter{
		emailAdapter: []noti.NotiAdapter{emailAdapter},
		pushAdapter:  []noti.NotiAdapter{pushAdapter},
	}

	for _, opt := range opts {
		opt(adapter)
	}

	return adapter
}

func (m *MultiNotiAdapter) Send(ctx context.Context, req model.NotiRequest) error {
	var targetAdapters []noti.NotiAdapter
	switch req.NotiType {
	case model.EmailType:
		targetAdapters = m.emailAdapter

	case model.PushType:
		targetAdapters = m.pushAdapter

	default:
		return errors.New("[MultiNotiAdapter.Send]: invalid notification type")
	}

	for _, adapter := range targetAdapters {
		if err := adapter.Send(ctx, req); err != nil {
			continue
		}

		return nil
	}

	return errors.New("[MultiNotiAdapter.Send]: failed to send notification on all adapters")
}
