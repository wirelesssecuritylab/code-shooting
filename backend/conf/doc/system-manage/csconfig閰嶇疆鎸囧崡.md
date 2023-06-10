# 部门和项目映射关系的配置
场景：
用户属于部门，部门属于研发中心，配置完成后，用户就可以查看（且只能查看）所关联研发中心的靶场。

配置方法：
在 conf/privilege 目录下，通过修改 projects.json 文件，对部门和研发中心的映射关系进行配置。

示例:
```json
    {
      "name": "软件研发中心",
      "project_id": "software-department",
      "mapping": {
        "departments": [
          "软件一部",
          "软件二部",
          "软件三部",
        ]
      }
    }
```

# 管理员的配置
场景：
用户默认具有打靶权限，管理员权限需要单独配置。

配置方法：
在 conf/privilege 目录下，修改 user-role-map.json 文件，对用户赋予或者移除管理员权限。

示例：
```json
        {
            "id": "123456789",
            "mapping": {
                "role_ids": [
                    "admin"
                ]
            }
        }
```
# 非本部门用户的配置
场景：
如果非本部门的用户要参加本部门的打靶，可以将用户手动配置给本部门。

配置方法：
在 conf/project 目录下，修改 deptUserMapping.json 文件，将非部门的用户添加到所属部门。

示例：
```json
    {
      "depart_name": "软件一部",
      "mapping": {
        "staff_id": [
          "8982293781"
        ]
      }
    }
```