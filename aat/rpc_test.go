package aat

import (
	"context"
	"testing"
	"time"

	"github.com/gammazero/nexus/client"
	"github.com/gammazero/nexus/wamp"
)

func TestRPCRegisterAndCall(t *testing.T) {
	// Connect callee session.
	callee, err := connectClient()
	if err != nil {
		t.Fatal("Failed to connect client: ", err)
	}

	// Test registering a valid procedure.
	handler := func(ctx context.Context, args []interface{}, kwargs, details map[string]interface{}) *client.InvokeResult {
		var sum int64
		for i := range args {
			n, ok := wamp.AsInt64(args[i])
			if ok {
				sum += n
			}
		}
		return &client.InvokeResult{Args: []interface{}{sum}}
	}

	// Register procedure "sum"
	procName := "sum"
	if err = callee.Register(procName, handler, nil); err != nil {
		t.Fatal("failed to register procedure: ", err)
	}

	// Connect caller session.
	caller, err := connectClient()
	if err != nil {
		t.Fatal("Failed to connect client: ", err)
	}

	// Test calling the procedure.
	callArgs := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ctx := context.Background()
	result, err := caller.Call(ctx, procName, nil, callArgs, nil, "")
	if err != nil {
		t.Fatal("failed to call procedure: ", err)
	}
	sum, ok := wamp.AsInt64(result.Arguments[0])
	if !ok {
		t.Fatal("could not convert result to int64")
	}
	if sum != 55 {
		t.Fatal("wrong result:", sum)
	}

	// Test unregister.
	if err = callee.Unregister(procName); err != nil {
		t.Fatal("failed to unregister procedure: ", err)
	}

	err = caller.Close()
	if err != nil {
		t.Fatal("Failed to disconnect client: ", err)
	}

	err = callee.Close()
	if err != nil {
		t.Fatal("Failed to disconnect client: ", err)
	}
}

func TestRPCCallUnregistered(t *testing.T) {
	// Connect caller session.
	caller, err := connectClient()
	if err != nil {
		t.Fatal("Failed to connect client: ", err)
	}

	// Test calling unregistered procedure.
	callArgs := []interface{}{555}
	ctx := context.Background()
	result, err := caller.Call(ctx, "NotRegistered", nil, callArgs, nil, "")
	if err == nil {
		t.Fatal("expected error calling unregistered procedure")
	}
	if result != nil {
		t.Fatal("result should be nil on error")
	}

	err = caller.Close()
	if err != nil {
		t.Fatal("Failed to disconnect client: ", err)
	}
}

func TestRPCUnregisterUnregistered(t *testing.T) {
	// Connect caller session.
	callee, err := connectClient()
	if err != nil {
		t.Fatal("Failed to connect client: ", err)
	}

	// Test unregister unregistered procedure.
	if err = callee.Unregister("NotHere"); err == nil {
		t.Fatal("expected error unregistering unregistered procedure")
	}

	err = callee.Close()
	if err != nil {
		t.Fatal("Failed to disconnect client: ", err)
	}
}

func TestRPCCancelCall(t *testing.T) {
	// Connect callee session.
	callee, err := connectClient()
	if err != nil {
		t.Fatal("Failed to connect client: ", err)
	}

	// Register procedure that waits.
	handler := func(ctx context.Context, args []interface{}, kwargs, details map[string]interface{}) *client.InvokeResult {
		<-ctx.Done() // handler will block forever until canceled.
		return &client.InvokeResult{Err: wamp.ErrCanceled}
	}
	procName := "myproc"
	if err = callee.Register(procName, handler, nil); err != nil {
		t.Fatal("failed to register procedure: ", err)
	}

	// Connect caller session.
	caller, err := connectClient()
	if err != nil {
		t.Fatal("Failed to connect client: ", err)
	}

	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	// Calling the procedure, should block.
	go func() {
		callArgs := []interface{}{73}
		_, err := caller.Call(ctx, procName, nil, callArgs, nil, "killnowait")
		errChan <- err
	}()

	// Make sure the call is blocked.
	select {
	case err = <-errChan:
		t.Fatal("call should have been blocked")
	case <-time.After(200 * time.Millisecond):
	}

	cancel()

	// Make sure the call is canceled.
	select {
	case err = <-errChan:
	case <-time.After(2 * time.Second):
		t.Fatal("call should have been canceled")
	}

	rpcError, ok := err.(client.RPCError)
	if !ok {
		t.Fatal("expected RPCError type of error")
	}
	if rpcError.Err.Error != wamp.ErrCanceled {
		t.Fatal("expected canceled error, got:", err)
	}
	if err = callee.Unregister(procName); err != nil {
		t.Fatal("failed to unregister procedure: ", err)
	}
	err = callee.Close()
	if err != nil {
		t.Fatal("Failed to disconnect client: ", err)
	}
	err = caller.Close()
	if err != nil {
		t.Fatal("Failed to disconnect client: ", err)
	}
}
