import colors from '@/constants/colors';
import React, { useRef, useState } from 'react';
import { StyleSheet, Text, TextInput, TextInputProps, View } from 'react-native';

type FormInputProps = {
  title: string;
  placeholder?: string;
  value: string;
  onChangeText: (text: string) => void;
  type?: 'default' | 'mail' | 'password';
};

const FormInput = ({ title, placeholder, value, onChangeText, type }: FormInputProps) => {
  const inputRef = useRef<TextInput>(null);
  const [secure, setSecure] = useState(false);

  const baseProps: TextInputProps = {
    style: styles.input,
    placeholder: placeholder,
    value: value,
    onChangeText: onChangeText,
  };
  let typeProps: TextInputProps = {};
  switch (type) {
    case 'mail':
      typeProps = {
        keyboardType: 'email-address',
        textContentType: 'emailAddress',
        autoComplete: 'email',
        importantForAutofill: 'yes',
        onTouchStart() {
          setSecure(false);
        },
      };
      break;

    case 'password':
      typeProps = {
        /*
          フォーカス時以外はパスワード入力モードをオフ
          サードパーティー製キーボードが開なるため応急処置
        */
        secureTextEntry: secure,
        textContentType: secure ? 'password' : 'none',
        autoComplete: 'password',
        importantForAutofill: 'yes',
        returnKeyType: 'done',
        onFocus: () => {
          setSecure(true);
          inputRef.current?.focus();
        },
        onBlur: () => {
          setSecure(false);
        },
      };
      break;
    default:
      typeProps = {
        returnKeyType: 'next',
        onTouchStart() {
          setSecure(false);
        },
      };
  }

  return (
    <View style={styles.container}>
      <Text style={styles.text}>{title}</Text>

      <TextInput
        {...baseProps}
        {...typeProps}
        ref={inputRef}
        onTouchStart={() => {
          inputRef.current?.blur();
          setTimeout(() => {
            inputRef.current?.focus();
          }, 100);
        }}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    alignItems: 'flex-start',
    alignSelf: 'stretch',
    gap: 8,
  },
  input: {
    alignSelf: 'stretch',
    backgroundColor: colors.light.card.input,
    borderColor: colors.light.card.border,
    borderRadius: 8,
    borderStyle: 'solid',
    borderWidth: 1,
    height: 36,
    paddingHorizontal: 13,
    paddingVertical: 9,
  },
  text: {
    color: colors.light.card.sectionTitle,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 20,
  },
});

export default FormInput;
