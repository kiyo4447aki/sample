//TODO: ログイン状態により、ナビゲーションを振り分けるように実装
//import { useAuth } from '../contexts/AuthContext';
import { useAuth } from '@/contexts/AuthContext';
import CameraDetailScreen from '@/screens/CameraDetailScreen/CameraDetailScreen';
import LiveViewScreen from '@/screens/LiveViewScreen/LiveViewScreen';
import LoginScreen from '@/screens/LoginScreen/LoginScreen';
import PlayRecordScreen from '@/screens/PlayRecordScreen/PlayRecordScreen';
import RecListScreen from '@/screens/RecListScreen/RecListScreen';
import TopScreen from '@/screens/TopScreen';
import { NavigationContainer } from '@react-navigation/native';
import { createStackNavigator } from '@react-navigation/stack';

export type RootStackParamList = {
  Top: undefined;
  RecList: undefined;
  CameraDetail: {
    cameraId: string;
  };
  LiveView: {
    cameraId: string;
  };
  // LiveViewAll: undefined;
  PlayRecord: { url: string };
};

const Stack = createStackNavigator<RootStackParamList>();

const RootNavigator = () => {
  //const context = useAuth();
  const { token } = useAuth();

  return token ? (
    <NavigationContainer>
      <Stack.Navigator
        screenOptions={{
          headerShown: false,
        }}
      >
        <Stack.Screen name="Top" component={TopScreen} />
        <Stack.Screen name="RecList" component={RecListScreen} />
        <Stack.Screen name="CameraDetail" component={CameraDetailScreen} />

        <Stack.Screen name="LiveView" component={LiveViewScreen} />
        {/*<Stack.Screen name="LiveViewAll" component={undefined} />*/}
        <Stack.Screen name="PlayRecord" component={PlayRecordScreen} />
      </Stack.Navigator>
    </NavigationContainer>
  ) : (
    <LoginScreen />
  );
};

export default RootNavigator;
