import React from 'react';
import {
  Image,
  ImageSourcePropType,
  ImageStyle,
  StyleSheet,
  Text,
  TextStyle,
  TouchableOpacity,
  View,
  ViewStyle,
} from 'react-native';

import arrowRight from '@/../assets/icons/caret-right.png';
import colors from '@/constants/colors';

type CardProps = {
  style?: ViewStyle;
  onPress?: () => void;
};

type CardComposition = {
  Icon: React.FC<{ style?: ImageStyle; source: ImageSourcePropType }>;
  Title: React.FC<{ style?: TextStyle; text: string }>;
  Subtitle: React.FC<{ style?: TextStyle; text: string }>;
  TextContainer: React.FC<{ style?: ViewStyle; children: React.ReactNode }>;
  Arrow: React.FC;
};

const Card: React.FC<CardProps & { children: React.ReactNode }> & CardComposition = ({
  style,
  onPress,
  children,
}) => {
  return (
    <TouchableOpacity style={[styles.card, style]} onPress={onPress}>
      {children}
    </TouchableOpacity>
  );
};

Card.Icon = ({ style, source }) => {
  return <Image source={source} style={[styles.icon, style]} />;
};
Card.Title = ({ style, text }) => {
  return <Text style={[styles.title, style]}> {text}</Text>;
};
Card.Subtitle = ({ style, text }) => {
  return <Text style={[styles.subtitle, style]}> {text}</Text>;
};

Card.TextContainer = ({ style, children }) => {
  return <View style={[styles.textContainer, style]}>{children}</View>;
};

Card.Arrow = () => {
  return <Image source={arrowRight} style={styles.arrow} />;
};

const styles = StyleSheet.create({
  arrow: {
    height: 20,
    width: 10,
  },
  card: {
    alignItems: 'center',
    alignSelf: 'stretch',
    backgroundColor: colors.light.card.bgSolid,
    borderColor: colors.light.card.border,
    borderRadius: 14,
    borderStyle: 'solid',
    borderWidth: 1,
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 20,
  },
  icon: {
    height: 30,
    width: 30,
  },
  subtitle: {
    color: colors.light.card.secondaryText,
    fontSize: 14,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 20,
  },
  textContainer: {
    alignItems: 'flex-start',
    flexDirection: 'column',
    justifyContent: 'center',
  },
  title: {
    color: colors.light.card.primaryText,
    fontSize: 16,
    fontStyle: 'normal',
    fontWeight: 300,
    lineHeight: 24,
  },
});

export default Card;
