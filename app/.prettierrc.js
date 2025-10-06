module.exports = {
  // 基本設定
  printWidth: 100, // 1行の最大文字数
  tabWidth: 2, // インデントのスペース数
  useTabs: false, // スペースの代わりにタブを使用しない
  semi: true, // 文末にセミコロンを付ける
  singleQuote: true, // シングルクォートを使用
  quoteProps: 'as-needed', // オブジェクトのプロパティ名のクォートは必要な場合のみ
  jsxSingleQuote: false, // JSXではダブルクォートを使用

  // 末尾のカンマ設定
  trailingComma: 'all', // 可能な限り末尾にカンマを付ける

  // スペース設定
  bracketSpacing: true, // オブジェクトリテラルの括弧の間にスペースを入れる
  arrowParens: 'avoid', // アロー関数の引数が1つの場合は括弧を省略

  // JSX設定
  jsxBracketSameLine: false, // JSXの閉じ括弧を新しい行に配置

  // その他
  endOfLine: 'lf', // 改行コードをLFに統一
  embeddedLanguageFormatting: 'auto', // 埋め込みコードのフォーマットを自動で処理
};
