import { ref } from 'vue'

export type Locale = 'en' | 'ja'

// localStorage に保存して再起動後も維持する
const saved = localStorage.getItem('locale') as Locale | null
const locale = ref<Locale>(saved === 'en' || saved === 'ja' ? saved : 'ja')

const messages: Record<Locale, Record<string, string>> = {
  en: {
    'app.subtitle': 'Prompt Version Control',
    'lang.toggle': 'JA',
    'nav.prompts': 'Prompts',
    'nav.help': 'Help',
    'nav.settings': 'Settings',

    'error.prefix': 'Error: ',
    'error.retry': 'Retry',
    'loading': 'Loading…',

    'prompts.title': 'Prompts',
    'prompts.empty': 'No prompts yet. Create one to start tracking its history.',
    'prompts.new': 'New prompt',
    'prompts.new.key': 'Key (unique id)',
    'prompts.new.description': 'Description',
    'prompts.new.content': 'Initial content',
    'prompts.new.create': 'Create',
    'prompts.card.open': 'Open',
    'prompts.card.delete': 'Delete',
    'prompts.card.delete.confirm': 'Delete this prompt and all its history? This cannot be undone.',

    'detail.back': 'Back to prompts',
    'detail.newVersion': 'New version',
    'detail.newVersion.content': 'Content',
    'detail.newVersion.message': 'Message',
    'detail.newVersion.save': 'Save as new version',
    'detail.history': 'History',
    'detail.history.current': 'current',
    'detail.history.rollback': 'Roll back to this',
    'detail.history.rollback.confirm': 'Roll back to this version? History is kept — this only moves the current pointer.',
    'detail.diff.title': 'Compare',
    'detail.diff.from': 'From',
    'detail.diff.to': 'To',
    'detail.diff.empty': 'Pick two versions to compare.',
    'detail.quality.title': 'Quality trend',
    'detail.quality.empty': 'No quality ratings recorded yet.',
    'detail.quality.rated': '{n} ratings',

    'help.title': 'Help',
    'help.intro': 'How versions, rollback, and quality tracking fit together.',
    'help.what.title': 'What this app does',
    'help.what.body': 'Every prompt or worker definition you manage here keeps a full, append-only history. Rolling back never deletes anything — it just moves which version is "current". Compare any two versions side by side, and if an external system posts quality scores per version, you can see exactly which change made things better or worse.',
    'help.start.title': 'Getting started',
    'help.start.1': 'Create a prompt — this saves version 1 with the content you give it.',
    'help.start.2': 'Edit and save again to create version 2, 3, and so on. Nothing is ever overwritten.',
    'help.start.3': 'Use Compare to see exactly what changed between any two versions.',
    'help.start.4': 'If quality dropped after a change, roll back to the earlier version — the history stays intact either way.',
    'help.stuck.title': 'Common snags',
    'help.stuck.1': 'Quality trend is empty → nothing has POSTed to /api/v1/ratings for this prompt\'s versions yet.',
    'help.stuck.2': 'Rollback doesn\'t seem to do anything to the list → that\'s expected. It only changes which version is marked "current"; the version list itself never shrinks.',
    'help.stuck.3': 'Want real A/B traffic splitting → that\'s out of scope here; this tool tracks history and quality, not live routing.',

    'settings.title': 'Settings',
    'settings.api.title': 'API endpoint',
    'settings.api.desc': 'External systems can read/write prompts and post quality ratings here.',
    'settings.keys.title': 'API keys',
    'settings.keys.name': 'Key name',
    'settings.keys.issue': 'Issue key',
    'settings.keys.issued': 'Key issued — copy it now, it will not be shown again',
    'settings.keys.copy': 'Copy',
    'settings.keys.revoke': 'Revoke',
    'settings.keys.revoked': 'Revoked',
    'settings.keys.empty': 'No keys issued yet.',
    'settings.version': 'Version',
    'settings.quit': 'Quit',
    'settings.quit.confirm': 'Quit the app?',
  },
  ja: {
    'app.subtitle': 'Prompt Version Control',
    'lang.toggle': 'EN',
    'nav.prompts': 'プロンプト',
    'nav.help': 'ヘルプ',
    'nav.settings': '設定',

    'error.prefix': 'エラー: ',
    'error.retry': '再試行',
    'loading': '読み込み中…',

    'prompts.title': 'プロンプト',
    'prompts.empty': 'まだプロンプトがありません。作成して履歴管理を始めましょう。',
    'prompts.new': '新規作成',
    'prompts.new.key': 'Key（一意なID）',
    'prompts.new.description': '説明',
    'prompts.new.content': '初期内容',
    'prompts.new.create': '作成',
    'prompts.card.open': '開く',
    'prompts.card.delete': '削除',
    'prompts.card.delete.confirm': 'このプロンプトと全履歴を削除しますか？元に戻せません。',

    'detail.back': 'プロンプト一覧へ戻る',
    'detail.newVersion': '新しいバージョン',
    'detail.newVersion.content': '内容',
    'detail.newVersion.message': 'メッセージ',
    'detail.newVersion.save': '新バージョンとして保存',
    'detail.history': '履歴',
    'detail.history.current': '現在地点',
    'detail.history.rollback': 'ここへロールバック',
    'detail.history.rollback.confirm': 'このバージョンへロールバックしますか？履歴は保持されます — 現在地点のポインタが動くだけです。',
    'detail.diff.title': '比較',
    'detail.diff.from': 'From',
    'detail.diff.to': 'To',
    'detail.diff.empty': '比較する2つのバージョンを選んでください。',
    'detail.quality.title': '品質推移',
    'detail.quality.empty': 'まだ品質スコアが記録されていません。',
    'detail.quality.rated': '{n}件の評価',

    'help.title': 'ヘルプ',
    'help.intro': 'バージョン管理・ロールバック・品質トラッキングがどう連動するかをまとめました。',
    'help.what.title': 'このアプリでできること',
    'help.what.body': 'ここで管理するプロンプト/worker定義は、完全な追記専用の履歴を持ちます。ロールバックは何も削除しません — 「現在」がどのバージョンかを動かすだけです。任意の2バージョンを並べて比較でき、外部システムがバージョンごとに品質スコアをPOSTすれば、どの変更で良くなった/悪くなったかが正確に分かります。',
    'help.start.title': 'はじめに',
    'help.start.1': 'プロンプトを作成する — 入力した内容がバージョン1として保存されます。',
    'help.start.2': '編集して再保存するとバージョン2、3…と増えていきます。上書きされることはありません。',
    'help.start.3': '「比較」で任意の2バージョン間の差分を正確に確認できます。',
    'help.start.4': '変更後に品質が下がったら、以前のバージョンへロールバック。どちらの場合も履歴はそのまま残ります。',
    'help.stuck.title': 'よくある詰まりどころ',
    'help.stuck.1': '品質推移が空 → このプロンプトのバージョンに対してまだ/api/v1/ratingsへのPOSTが届いていません。',
    'help.stuck.2': 'ロールバックしても一覧が変わらない → それが正しい挙動です。「現在」のマークが動くだけで、バージョン一覧自体は減りません。',
    'help.stuck.3': '実トラフィックのA/B振り分けをしたい → 本製品の対象外です。履歴と品質の追跡に特化しており、実際のルーティングは行いません。',

    'settings.title': '設定',
    'settings.api.title': 'APIエンドポイント',
    'settings.api.desc': '外部システムはここでプロンプトの読み書き・品質スコアの送信ができます。',
    'settings.keys.title': 'APIキー',
    'settings.keys.name': 'キー名',
    'settings.keys.issue': 'キーを発行',
    'settings.keys.issued': 'キーを発行しました — この場では二度と表示されないので今すぐコピーしてください',
    'settings.keys.copy': 'コピー',
    'settings.keys.revoke': '失効',
    'settings.keys.revoked': '失効済み',
    'settings.keys.empty': 'まだキーがありません。',
    'settings.version': 'バージョン',
    'settings.quit': '終了',
    'settings.quit.confirm': 'アプリを終了しますか？',
  },
}

export function useI18n() {
  function t(key: string, params?: Record<string, string | number>): string {
    let msg = messages[locale.value][key] ?? messages.en[key] ?? key
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        msg = msg.replace(`{${k}}`, String(v))
      }
    }
    return msg
  }

  function toggleLocale() {
    locale.value = locale.value === 'en' ? 'ja' : 'en'
    localStorage.setItem('locale', locale.value)
  }

  return { t, locale, toggleLocale }
}
