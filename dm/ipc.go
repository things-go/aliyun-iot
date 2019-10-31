package dm

type ipcEvtType byte

const (
	ipcEvtRawReply ipcEvtType = iota
)

type ipcMessage struct {
	err        error
	evt        ipcEvtType
	devID      int
	productKey string
	deviceName string
	payload    interface{}
}

func (sf *Client) ipcSendMessage(msg *ipcMessage) error {
	select {
	case sf.ipc <- msg:
	default:
		return ErrIPCMessageBuffFull
	}
	return nil
}

func (sf *Client) ipcRunMessage() {
	for {
		select {
		case msg := <-sf.ipc:
			switch msg.evt {

			}
		}
	}
}
