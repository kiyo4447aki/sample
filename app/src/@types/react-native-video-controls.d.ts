//最小限の型定義のみ行うため、any型を許容する
//TODO anyの削除と型定義
/* eslint-disable @typescript-eslint/no-explicit-any */

// VideoPlayer.d.ts
declare module 'react-native-video-controls' {
  import { Component } from 'react';
  import { StyleProp, ViewStyle } from 'react-native';
  import { ReactVideoProps, ResizeMode } from 'react-native-video';

  export interface VideoPlayerProps extends ReactVideoProps {
    /**
     * フルスクリーン切り替え時に resizeMode を自動で切り替えるか
     * @default true
     */
    toggleResizeModeOnFullscreen?: boolean;
    /**
     * コントロール表示／非表示アニメーションの所要ミリ秒
     * @default 500
     */
    controlAnimationTiming?: number;
    /**
     * ダブルタップと判定するまでの最大ミリ秒
     * @default 130
     */
    doubleTapTime?: number;
    /**
     * バックグラウンド再生を許可するか
     * @default false
     */
    playInBackground?: boolean;
    /**
     * Inactive 時再生を継続するか
     * @default false
     */
    playWhenInactive?: boolean;
    /**
     * デフォルトのリサイズモード
     * @default 'contain'
     */
    resizeMode?: ResizeMode;
    /**
     * 初期表示時にフルスクリーンモードにするか
     * @default false
     */
    isFullscreen?: boolean;
    /**
     * マウント時にコントロールを表示しておくか
     * @default true
     */
    showOnStart?: boolean;
    /**
     * 初期状態の再生／一時停止
     * @default false
     */
    paused?: boolean;
    /**
     * ループ再生するか
     * @default false
     */
    repeat?: boolean;
    /**
     * 初期ミュート状態
     * @default false
     */
    muted?: boolean;
    /**
     * 初期音量 (0〜1)
     * @default 1
     */
    volume?: number;
    /**
     * タイトル文字列
     * @default ''
     */
    title?: string;
    // ページバック時に利用するナビゲーター
    navigator?: {
      pop: () => void;
    };
    /**
     * 再生速度
     * @default 1
     */
    rate?: number;
    /**
     * コントロールを自動で隠すまでの遅延 (ms)
     * @default 15000
     */
    controlTimeout?: number;
    /**
     * シークバー中の scrubbing 間隔 (ms)
     * @default 0
     */
    scrubbing?: number;
    /**
     * 画面タップで再生／一時停止を許可するか
     */
    tapAnywhereToPause?: boolean;

    // シークバーのカラー
    seekColor?: string;

    /** 全体のコンテナスタイル */
    style?: StyleProp<ViewStyle>;
    /** <Video> 部分のスタイル */
    videoStyle?: StyleProp<ViewStyle>;

    /** コントロールを無効にするフラグ群 */
    disableBack?: boolean;
    disableVolume?: boolean;
    disableFullscreen?: boolean;
    disableTimer?: boolean;
    disableSeekbar?: boolean;

    /** 各種イベントコールバック */
    onError?: (error: any) => void;
    onBack?: () => void;
    onEnd?: () => void;
    onEnterFullscreen?: () => void;
    onExitFullscreen?: () => void;
    onShowControls?: () => void;
    onHideControls?: () => void;
    onLoadStart?: (data: any) => void;
    onProgress?: (data: any) => void;
    onSeek?: (data: any) => void;
    onLoad?: (data: any) => void;
    onPause?: () => void;
    onPlay?: () => void;
  }

  export default class VideoPlayer extends Component<VideoPlayerProps> {}
}
