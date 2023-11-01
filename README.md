
[TOC]

# nunu

- 快速產生CRUD template的工具

## tree exmpale
```
└─template
    └─nunu
      ├─append.yaml
      ├─handler.tpl
      ├─model_bo.tpl
      ├─model_dto.tpl
      ├─model_po.tpl
      ├─repository.tpl
      ├─service.tpl
      └─target.yaml
```

- target.yaml example

```yaml
target_path:
  handler: "./service/internal/controller/handler/"
  service: "./service/internal/core/module/"
  repository: "./service/internal/repository/"
  model: 
    - path: "./service/internal/model/bo/"
      temp_file: "model_bo.tpl"
    - path: "./service/internal/model/po/"
      temp_file: "model_po.tpl"
    - path: "./service/internal/model/dto/"
      temp_file: "model_dto.tpl"
```

- append.yaml example

```yaml
files:
  - path: "internal/handler/provider.go"
    regex: "type Handler struct \\{.*?\\}"
    template: "\\n\\t{{ .FileName }}Handler  *{{ .FileNameTitleLower }}Handler"
  - path: "internal/handler/provider.go"
    regex: "func NewHandler\\(in digIn\\) .*?Handler \\{.*?\\\t}"
    template: "\\n\\t\\t{{ .FileName }}Handler:  &{{ .FileNameTitleLower }}Handler{in: in},"
  - path: "internal/service/provider.go"
    regex: "type Service struct \\{.*?\\}"
    template: "\\n\\t{{ .FileName }}Srv  I{{ .FileName }}Service"
  - path: "internal/service/provider.go"
    regex: "func NewService\\(in digIn\\) .*?Service \\{.*?\\\t}"
    template: "\\n\\t\\t{{ .FileName }}Srv:  New{{ .FileName }}Service(in),"
  - path: "internal/repository/provider.go"
    regex: "type Repository struct \\{.*?\\}"
    template: "\\n\\t{{ .FileName }}Repo  I{{ .FileName }}Repository"
  - path: "internal/repository/provider.go"
    regex: "func NewRepository\\(in digIn\\) .*?Repository \\{.*?\\\t}"
    template: "\\n\\t\\t{{ .FileName }}Repo:  New{{ .FileName }}Repository(in),"
```


## 指令(command)

### 創建(create)

- nunu create handler `{model}`
  - auto gen handler.tpl to target directory
- nunu create service `{model}`
  - auto gen handler.tpl to target directory
- nunu create reposotory `{model}`
  - auto gen handler.tpl to target directory
- nunu create model `{model}`
  - auto gen handler.tpl to target directory
- nunu create all `{model}`
  - auto gen handler & service & reposotory & model

### 插入append.yaml內容行代碼(append)

- nunu append `{model}`
  - 主要藉由正則以及模板，達成在需要依賴註冊的地方添加自己需要的代碼

### 套件升級(upgrade)

- nunu upgrade

### 其他

- 其餘功能同fork的來源倉儲nunu
