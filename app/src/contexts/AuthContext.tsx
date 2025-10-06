//TODO:ログインのリトライ処理
//TODO:トークンの更新処理
import {
  createContext,
  PropsWithChildren,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react';
import api from '../services/api';
import { loadToken, saveToken } from '../services/token';

type LoginResponse = {
  status: 'success' | 'error';
  token?: string;
  error?: string;
};

type AuthContextType = {
  token: string;
  signIn: (id: string, password: string, isKeepLogin: boolean) => Promise<void>;
  signOut: () => Promise<void>;
  isLoggedIn: boolean;
  error: string;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const useAuth = () => {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuthはAuthProviderの中で使用してください');
  return ctx;
};

const AuthProvider = ({ children }: PropsWithChildren) => {
  const [token, setToken] = useState<string>('');
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    loadToken()
      //TODO 自動でトークンリフレッシュ
      .then(t => setToken(t))
      .catch((e: Error) => {
        setError(e.toString());
        setToken('');
      })
      .finally(() => {});
  }, []);

  useEffect(() => {
    if (token) {
      api.defaults.headers.common.Authorization = `Bearer ${token}`;
      setIsLoggedIn(true);
    } else {
      delete api.defaults.headers.common.Authorization;
    }
  }, [token]);

  const signIn = useCallback(async (username: string, password: string, isKeepLogin: boolean) => {
    setIsLoggedIn(false);
    setError('');

    try {
      const res = await api.post<LoginResponse>('/login', {
        username,
        password,
      });
      if (res.data.status !== 'success' || typeof res.data.token !== 'string') {
        throw new Error(`認証に失敗しました エラー詳細：\n${res.data.error}`);
      }
      setToken(res.data.token);
      if (isKeepLogin) {
        await saveToken(res.data.token);
      }
      setIsLoggedIn(true);
    } catch (e: unknown) {
      if (e instanceof Error) {
        setError(e.message);
      }
      throw e;
    }
  }, []);

  const signOut = useCallback(async () => {
    setToken('');
    await saveToken('');
  }, []);

  const value = useMemo(
    () => ({
      token,
      signIn,
      signOut,
      isLoggedIn,
      error,
    }),
    [token, signIn, signOut, isLoggedIn, error],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export { AuthProvider, useAuth };
