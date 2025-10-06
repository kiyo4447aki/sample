//TODO: エラーコードの実装
//TODO: 再接続処理の実装
import { useCallback, useEffect, useRef, useState } from 'react';
import { RTCIceCandidate, RTCPeerConnection, RTCSessionDescription } from 'react-native-webrtc';

type UseJanusParams = {
  serverURL?: string;
  roomID: number;
  feedID: number;
  pass?: number;
};

const useJanus = ({ serverURL = '', roomID, feedID, pass }: UseJanusParams) => {
  const [streamURL, setStreamURL] = useState<string>('');
  const [connecting, setConnecting] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  const wsRef = useRef<WebSocket | null>(null);
  const pcRef = useRef<RTCPeerConnection | null>(null);
  const txnRef = useRef<number>(0);
  const sessionIDRef = useRef<number>(0);
  const handleIDRef = useRef<number>(0);

  const makeTxn = () => `txn-${++txnRef.current}`;

  const disconnect = useCallback(() => {
    pcRef.current?.close();
    pcRef.current = null;
    wsRef.current?.close();
    wsRef.current = null;
    sessionIDRef.current = 0;
    handleIDRef.current = 0;
    txnRef.current = 0;
    setConnecting(false);
  }, []);

  useEffect(() => {
    if (!serverURL) return;
    setConnecting(true);
    setError('');
    setStreamURL('');

    const ws = new WebSocket(serverURL);
    wsRef.current = ws;

    //Janusセッションの作成
    ws.onopen = () => {
      ws.send(
        JSON.stringify({
          janus: 'create',
          transaction: makeTxn(),
        }),
      );
    };

    ws.onerror = () => {
      setError('通信エラーが発生しました');
      setConnecting(false);
      disconnect();
    };

    ws.onmessage = async ({ data }) => {
      //TODO メッセージの型定義
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      let msg: any;

      try {
        msg = JSON.parse(data);
      } catch {
        setError('シグナリング時に無効なメッセージを受け取りました');
        disconnect();
        return;
      }

      //セッション作成の成功時
      if (msg.janus === 'success' && msg.data?.id && !sessionIDRef.current) {
        sessionIDRef.current = msg.data.id;
        ws.send(
          JSON.stringify({
            janus: 'attach',
            plugin: 'janus.plugin.videoroom',
            session_id: sessionIDRef.current,
            transaction: makeTxn(),
          }),
        );
        //プラグインのアタッチ成功時
      } else if (
        msg.janus === 'success' &&
        msg.data?.id &&
        sessionIDRef.current &&
        !handleIDRef.current
      ) {
        handleIDRef.current = msg.data.id;
        ws.send(
          JSON.stringify({
            janus: 'message',
            session_id: sessionIDRef.current,
            handle_id: handleIDRef.current,
            body: {
              request: 'join',
              room: roomID,
              ptype: 'subscriber',
              feed: feedID,
              pin: pass,
            },
            transaction: makeTxn(),
          }),
        );
        //offer受信時
      } else if (msg.janus === 'event' && msg.jsep?.type === 'offer') {
        const pc = new RTCPeerConnection({
          iceServers: [{ urls: 'stun:stun.l.google.com:19302' }],
        });
        pcRef.current = pc;

        pc.addTransceiver('video', { direction: 'recvonly' });

        //トラック受信時にメディアストリームを取得
        pc.addEventListener('track', e => {
          if (e.streams && e.streams?.length > 0) {
            setStreamURL(e.streams[0].toURL());
            setConnecting(false);
          }
        });

        try {
          await pc.setRemoteDescription(new RTCSessionDescription(msg.jsep));
        } catch {
          setError('set offer failed');
          setConnecting(false);
          return;
        }
        let answer;
        try {
          answer = await pc.createAnswer();
        } catch {
          setError('create answer failed');
          setConnecting(false);
          return;
        }

        try {
          await pc.setLocalDescription(answer);
        } catch {
          setError('set answer failed');
          setConnecting(false);
          return;
        }

        //startメッセージ&local sdpを送信
        ws.send(
          JSON.stringify({
            janus: 'message',
            session_id: sessionIDRef.current,
            handle_id: handleIDRef.current,
            transaction: makeTxn(),
            body: { request: 'start', room: roomID },
            jsep: answer,
          }),
        );

        //ice candidates送信
        pc.addEventListener('icecandidate', e => {
          ws.send(
            JSON.stringify({
              janus: 'trickle',
              session_id: sessionIDRef.current,
              handle_id: handleIDRef.current,
              transaction: makeTxn(),
              candidate: e.candidate ? e.candidate : { completed: true },
            }),
          );
        });
      } else if (msg.janus === 'trickle' && msg.candidate) {
        try {
          await pcRef.current?.addIceCandidate(new RTCIceCandidate(msg.candidate));
        } catch {
          setError('ice candidatesの追加に失敗しました');
          setConnecting(false);
        }
      } else if (msg.janus === 'error') {
        setError(msg.error?.reason || 'Janus エラー');
        setConnecting(false);
        disconnect();
      }
    };

    ws.onclose = () => {
      setConnecting(false);
    };

    return () => {
      disconnect();
    };
  }, [disconnect]);

  return { streamURL, connecting, error, disconnect };
};

export default useJanus;
