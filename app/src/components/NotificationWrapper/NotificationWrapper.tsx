import { useAuth } from '@/contexts/AuthContext';
import { initPush } from '@/services/notification';
import React, { PropsWithChildren, useEffect, useRef } from 'react';
import { Alert } from 'react-native';

const NotificationWrapper = ({ children }: PropsWithChildren) => {
  const { isLoggedIn } = useAuth();
  const unsubRef = useRef<(() => void) | undefined>(undefined);

  useEffect(() => {
    if (!isLoggedIn) {
      if (unsubRef.current) {
        unsubRef.current();
        unsubRef.current = undefined;
      }
      return;
    }

    let cancelled = false;

    (async () => {
      try {
        const unsub = await initPush();
        if (cancelled) {
          if (typeof unsub === 'function') unsub();
          return;
        }
        unsubRef.current = typeof unsub === 'function' ? unsub : undefined;
      } catch {
        Alert.alert('通知の初期化に失敗しました', 'アプリを再起動してください');
      }
    })();

    return () => {
      cancelled = true;
      if (unsubRef.current) {
        unsubRef.current();
        unsubRef.current = undefined;
      }
    };
  }, [isLoggedIn]);

  return <>{children}</>;
};

export default NotificationWrapper;
