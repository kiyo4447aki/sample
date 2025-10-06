import lockIcon from '@/../assets/icons/lock.png';
import colors from '@/constants/colors';
import { useAuth } from '@/contexts/AuthContext';
import Button from '@/screens/LoginScreen/components/LoginForm/Button';
import FormInput from '@/screens/LoginScreen/components/LoginForm/FormInput';
import Checkbox from 'expo-checkbox';
import React, { useState } from 'react';
import { Alert, StyleSheet, Text, View } from 'react-native';

const LoginForm = () => {
  const [mailAddress, setMailAddress] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [isKeepLogin, setIsKeepLogin] = useState<boolean>(true);

  const { signIn } = useAuth();

  return (
    <View style={styles.container}>
      <Text style={styles.title}>ログイン</Text>
      <FormInput
        title="メールアドレス"
        placeholder="name@example.com"
        value={mailAddress}
        onChangeText={text => setMailAddress(text)}
        type="mail"
      />
      <FormInput
        title="パスワード"
        placeholder="パスワードを入力"
        value={password}
        onChangeText={text => setPassword(text)}
        type="password"
      />
      <View style={styles.keepLoginContainer}>
        <Checkbox
          value={isKeepLogin}
          onValueChange={value => setIsKeepLogin(value)}
          style={styles.checkbox}
        />
        <Text style={styles.checkboxText}>ログイン状態を保持</Text>
      </View>
      <Button
        text="ログイン"
        icon={lockIcon}
        onPress={() => {
          signIn(mailAddress, password, isKeepLogin).catch(() =>
            Alert.alert('ログインに失敗しました', 'メールアドレスとパスワードをご確認ください'),
          );
        }}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  checkbox: {
    backgroundColor: colors.light.checkbox.bg,
    borderColor: colors.light.checkbox.border,
    borderRadius: 4,
    borderStyle: 'solid',
    borderWidth: 1,
    height: 16,
    width: 16,
  },
  checkboxText: {
    color: colors.light.checkbox.text,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 20,
  },
  container: {
    alignItems: 'flex-start',
    alignSelf: 'stretch',
    backgroundColor: colors.light.card.bg,
    borderColor: colors.light.card.border,
    borderRadius: 14,
    borderStyle: 'solid',
    borderWidth: 1,
    gap: 16,
    paddingHorizontal: 24,
    paddingVertical: 25,
  },
  keepLoginContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: 8,
  },
  title: {
    alignSelf: 'center',
    color: colors.light.card.primaryText,
    fontSize: 20,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 28,
    marginBottom: 8,
    textAlign: 'center',
  },
});

export default LoginForm;
