package main

const gitChgLogConf = `style: github
template: CHANGELOG.tpl.md
info:
  title: CHANGELOG
  repository_url: ${VCS_URL}
options:
  commits:
    filters:
      Type:
        - add
        - fix
        - dep
        - chg
        - rmv
  commit_groups:
    title_maps:
      add: Added
      fix: Fixed
      dep: Deprecated
      chg: Changed
      rmv: Removed
  header:
    pattern: "^(\\w*)\\:\\s(.*)$"
    pattern_maps:
      - Type
      - Subject
  notes:
    keywords:
      - BREAKING CHANGE
`
