import { JanusSignalingClient } from '@/lib/JanusSignalingClient';
import React, { useEffect, useRef, useState } from 'react';
import { MediaStream, RTCView } from 'react-native-webrtc';

type JanusViewerProps = {
  url: string;
  roomId: string;
  feedId: string;
  password: string;
};

const JanusViewer = ({ url, roomId, feedId, password }: JanusViewerProps) => {
  const janusClient = useRef<JanusSignalingClient | undefined>(undefined);
  const [remoteStream, setRemoteStream] = useState<MediaStream | undefined>(undefined);

  useEffect(() => {
    janusClient.current = new JanusSignalingClient({
      roomId: roomId,
      feedId: feedId,
      janusUrl: url,
      password: password,
      onRemoteStream: (stream: MediaStream) => setRemoteStream(stream), // 映像受信時の処理
    });
    janusClient.current.connect(); // 接続開始
    return () => janusClient.current?.disconnect(); // クリーンアップ処理
  }, []);

  return (
    <>
      {remoteStream ? (
        //<RTCView streamURL={remoteStream.toURL()} style={{ flex: 1 }} objectFit="contain" />
        <RTCView streamURL={remoteStream.toURL()} style={{ flex: 1 }} objectFit="contain" />
      ) : undefined}
    </>
  );
};

export default JanusViewer;
