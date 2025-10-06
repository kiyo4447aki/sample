//TODO ルームのパスワード利用

import {
  MediaStream,
  RTCIceCandidate,
  RTCPeerConnection,
  RTCSessionDescription,
} from 'react-native-webrtc';

export interface JanusSignalingOptions {
  janusUrl: string;
  roomId: string; // Janus のルームID
  password: string;
  feedId: string; // 受信対象の publisher feed ID
  onRemoteStream: (stream: MediaStream) => void; // 映像受信時のコールバック
}

export class JanusSignalingClient {
  private ws: WebSocket | null = null; // WebSocket 接続
  private pc: RTCPeerConnection | null = null; // WebRTC ピア接続
  private sessionId?: number; // Janus セッション ID
  private handleId?: number; // Janus ハンドル ID（プラグイン接続）

  private keepAliveTimer?: ReturnType<typeof setInterval>; //keepAliveのsetIntervalのIDを保持

  constructor(private options: JanusSignalingOptions) {}

  // ランダムな英数字トランザクション ID を生成
  private makeTxn(): string {
    return `txn-${Math.random().toString(36).substring(2, 14)}`;
  }

  // WebSocket を通して Janus にメッセージ送信
  private send(msg: any) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(msg));
    }
  }

  public connect(): void {
    if (!this.options.janusUrl) {
      return;
    }

    this.ws = new WebSocket(this.options.janusUrl, 'janus-protocol');

    this.ws.onopen = () => {
      this.keepAliveTimer = setInterval(() => {
        if (this.sessionId) {
          this.send({
            janus: 'keepalive',
            session_id: this.sessionId,
            transaction: this.makeTxn(),
          });
        }
      }, 5000);

      this.send({ janus: 'create', transaction: this.makeTxn() });
    };

    this.ws.onmessage = async (event: MessageEvent) => {
      const data = JSON.parse(event.data);

      // セッション作成応答
      if (data.janus === 'success' && data.data && data.data.id && !this.sessionId) {
        this.sessionId = data.data.id;
        this.send({
          janus: 'attach',
          plugin: 'janus.plugin.videoroom',
          session_id: this.sessionId,
          transaction: this.makeTxn(),
        });
      }
      // プラグインアタッチ応答
      else if (
        data.janus === 'success' &&
        data.data &&
        data.data.id &&
        !this.handleId &&
        this.sessionId
      ) {
        this.handleId = data.data.id;
        this.send({
          janus: 'message',
          session_id: this.sessionId,
          handle_id: this.handleId,
          body: {
            request: 'join',
            room: this.options.roomId,
            ptype: 'subscriber',
            feed: this.options.feedId,
          },
          transaction: this.makeTxn(),
        });
      }

      // SDP offer を受信した場合
      else if (data.janus === 'event' && data.jsep && data.jsep.type === 'offer') {
        this.pc = new RTCPeerConnection();
        // 受信用トランシーバーを明示的に追加
        this.pc.addTransceiver('video', { direction: 'recvonly' });

        this.pc.addEventListener('track', e => {
          if (e.streams && e.streams.length > 0) {
            this.options.onRemoteStream(e.streams[0]);
          }
        });

        try {
          await this.pc.setRemoteDescription(new RTCSessionDescription(data.jsep));
          const answer = await this.pc.createAnswer();
          await this.pc.setLocalDescription(answer);
          this.send({
            janus: 'message',
            session_id: this.sessionId,
            handle_id: this.handleId,
            body: { request: 'start', room: this.options.roomId },
            jsep: answer,
            transaction: this.makeTxn(),
          });
        } catch (err) {
          if (err instanceof Error) {
            throw new Error(`SDPの交換に失敗しました：${err}`);
          } else {
            throw new Error('SDPの交換中に不明なエラーが発生しました');
          }
        }

        // ICE candidate の処理（trickle ICE）
        this.pc.addEventListener('icecandidate', e => {
          if (e.candidate) {
            this.send({
              janus: 'trickle',
              session_id: this.sessionId,
              handle_id: this.handleId,
              transaction: this.makeTxn(),
              candidate: e.candidate,
            });
          } else {
            this.send({
              janus: 'trickle',
              session_id: this.sessionId,
              handle_id: this.handleId,
              transaction: this.makeTxn(),
              candidate: { completed: true },
            });
          }
        });
      }
      // trickle ICE の候補が届いた場合
      else if (data.janus === 'trickle' && data.candidate) {
        try {
          const candidate = new RTCIceCandidate(data.candidate);
          await this.pc?.addIceCandidate(candidate);
        } catch (err) {
          if (err instanceof Error) {
            throw new Error(`remote candidateを追加できませんでした：${err}`);
          } else {
            throw new Error('remote candidateの追加中に不明なエラーが発生しました');
          }
        }
      }
    };

    this.ws.onerror = () => {
      if (this.keepAliveTimer) {
        clearInterval(this.keepAliveTimer);
        this.keepAliveTimer = undefined;
      }
    };

    this.ws.onclose = () => {
      if (this.keepAliveTimer) {
        clearInterval(this.keepAliveTimer);
        this.keepAliveTimer = undefined;
      }
    };
  }

  public disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    if (this.pc) {
      this.pc.close();
      this.pc = null;
    }
    if (this.keepAliveTimer) {
      clearInterval(this.keepAliveTimer);
      this.keepAliveTimer = undefined;
    }
  }
}
