# hubi
hubi helps you win at GitHub Issues.  https://twitter.com/dddaisuke

GitHubのissueのタイトルと本文を他のリポジトリにコピーするコマンドです。

# configure
https://github.com/settings/applications#personal-access-tokens の`[Generate new token]ボタン`からアクセストークンを生成します。生成したアクセストークンを以下の位置にコピーします。

`~/.github/config.json`
```
{
  "AccessToken": "YOUR ACCESS TOKEN"
}
```

# usage
> $ icp issue番号 ターゲットとなるリポジトリ名 [-c]

`-c`オプションを付けると、コピー元のissueを`CLOSE`します。
