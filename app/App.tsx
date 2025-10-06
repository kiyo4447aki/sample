import NotificationWrapper from '@/components/NotificationWrapper/NotificationWrapper';
import { AuthProvider } from '@/contexts/AuthContext';
import { DevicesProvider } from '@/contexts/DevicesContext';
import RootNavigator from '@/navigation/RootNavigator';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { SafeAreaProvider } from 'react-native-safe-area-context';

//デバッグ用
(window as any).AsyncStorage = AsyncStorage;

export default function App() {
  return (
    <AuthProvider>
      <DevicesProvider>
        <NotificationWrapper>
          <SafeAreaProvider>
            <RootNavigator />
          </SafeAreaProvider>
        </NotificationWrapper>
      </DevicesProvider>
    </AuthProvider>
  );
}
