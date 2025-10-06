import logo from '@/../assets/images/logo.png';
import colors from '@/constants/colors';
import LoginForm from '@/screens/LoginScreen/components/LoginForm/LoginForm';
import React from 'react';
import { Dimensions, Image, Keyboard, StyleSheet, Text, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

const LoginScreen = () => {
  return (
    <SafeAreaView style={styles.container} onTouchStart={() => Keyboard.dismiss()}>
      <View style={styles.contentArea}>
        <Image style={styles.logo} source={logo} />
        <Text style={styles.heading}>防犯カメラシステムにログイン</Text>
        <LoginForm />
      </View>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    backgroundColor: colors.light.bg,
    flex: 1,
    justifyContent: 'flex-start',
  },
  contentArea: {
    alignItems: 'center',
    alignSelf: 'stretch',
    flex: 1,
    gap: 15,
    justifyContent: 'flex-start',
    paddingHorizontal: 15,
    paddingVertical: 10,
  },
  heading: {
    alignSelf: 'center',
    color: colors.light.card.secondaryText,
    fontSize: 16,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 24,
    textAlign: 'center',
  },
  logo: {
    height: (8.9 / 100) * Dimensions.get('window')['height'],
    width: (16.6 / 100) * Dimensions.get('window')['width'],
  },
});

export default LoginScreen;
