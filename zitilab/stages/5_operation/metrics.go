package zitilab_5_operation

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/netfoundry/fablab/kernel"
	"github.com/netfoundry/ziti-fabric/pb/mgmt_pb"
	"github.com/netfoundry/ziti-foundation/channel2"
	"github.com/netfoundry/ziti-foundation/identity/dotziti"
	"github.com/netfoundry/ziti-foundation/transport"
	"github.com/sirupsen/logrus"
	"time"
)

func Metrics(closer chan struct{}) kernel.OperatingStage {
	return &metrics{closer: closer}
}

func (metrics *metrics) Operate(m *kernel.Model) error {
	if endpoint, id, err := dotziti.LoadIdentity("fablab"); err == nil {
		if address, err := transport.ParseAddress(endpoint); err == nil {
			dialer := channel2.NewClassicDialer(id, address, nil)
			if ch, err := channel2.NewChannel("metrics", dialer, nil); err == nil {
				metrics.ch = ch
			} else {
				return fmt.Errorf("error connecting metrics channel (%w)", err)
			}
		} else {
			return fmt.Errorf("invalid endpoint address (%w)", err)
		}
	} else {
		return fmt.Errorf("unable to load 'fablab' identity (%w)", err)
	}

	metrics.ch.AddReceiveHandler(metrics)

	request := &mgmt_pb.StreamMetricsRequest{
		Matchers: []*mgmt_pb.StreamMetricsRequest_MetricMatcher{
			&mgmt_pb.StreamMetricsRequest_MetricMatcher{},
		},
	}
	body, err := proto.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling metrics request (%w)", err)
	}

	requestMsg := channel2.NewMessage(int32(mgmt_pb.ContentType_StreamMetricsRequestType), body)
	errCh, err := metrics.ch.SendAndSync(requestMsg)
	if err != nil {
		logrus.Fatalf("error queuing metrics request (%w)", err)
	}
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error sending metrics request (%w)", err)
		}

	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout")
	}

	go metrics.runMetrics(m)

	return nil
}

func (metrics *metrics) ContentType() int32 {
	return int32(mgmt_pb.ContentType_StreamMetricsEventType)
}

func (metrics *metrics) HandleReceive(msg *channel2.Message, _ channel2.Channel) {
	response := &mgmt_pb.StreamMetricsEvent{}
	err := proto.Unmarshal(msg.Body, response)
	if err != nil {
		logrus.Error("error handling metrics receive (%w)", err)
	}

	logrus.Infof("<= [%s]", response.SourceId)
}

func (metrics *metrics) runMetrics(m *kernel.Model) {
	logrus.Infof("starting")
	defer logrus.Infof("exiting")

	for {
		select {
		case <-metrics.closer:
			return
		}
	}
}

type metrics struct {
	ch     channel2.Channel
	closer chan struct{}
}