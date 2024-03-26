#### 1. 项目简介
简单的将 suno.ai web转为api接口

#### 2. 部署方式

将 config.yaml.example 重命名为 config.yaml 并修改其中的配置
```yaml
Server:
    Port: 3560
App:
    Client: #登录 suno.ai 后的 cookie中的__client=xxxxx 的值
```
启动服务
```shell
./sunoweb2api
```
#### 3. 使用方式
创建任务
```shell
curl --location --request POST 'localhost:3560/v2/generate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "gpt_description_prompt": "an atmospheric metal song about dancing all night long",
    "mv": "chirp-v3-0",
    "prompt": "",
    "make_instrumental": false
}'
```
查询任务
```shell
curl --location --request POST 'localhost:3560/v2/feed' \
--header 'Content-Type: application/json' \
--data-raw '{
    "ids":"id1,id2"
}'
```
其中id1,id2为创建任务返回json中的clips中的id，如
```json
{
  "id": "b577cbad-18c3-49bb-a3aa-e39f9990a9c2",
  "clips": [
    {
      "id": "id1",
      "video_url": "",
      "audio_url": "",
      "image_url": null,
      "image_large_url": null,
      "major_model_version": "v3",
      "model_name": "chirp-v3",
      "metadata": {
        "tags": null,
        "prompt": "",
        "gpt_description_prompt": "an atmospheric metal song about dancing all night long",
        "audio_prompt_id": null,
        "history": null,
        "concat_history": null,
        "type": "gen",
        "duration": null,
        "refund_credits": null,
        "stream": true,
        "error_type": null,
        "error_message": null
      },
      "is_liked": false,
      "user_id": "284ac3ca-a2cf-4b0c-a1b5-64d7315fbe28",
      "is_trashed": false,
      "reaction": null,
      "created_at": "2024-03-26T07:01:44.235Z",
      "status": "submitted",
      "title": "",
      "play_count": 0,
      "upvote_count": 0,
      "is_public": false
    },
    {
      "id": "id2",
      "video_url": "",
      "audio_url": "",
      "image_url": null,
      "image_large_url": null,
      "major_model_version": "v3",
      "model_name": "chirp-v3",
      "metadata": {
        "tags": null,
        "prompt": "",
        "gpt_description_prompt": "an atmospheric metal song about dancing all night long",
        "audio_prompt_id": null,
        "history": null,
        "concat_history": null,
        "type": "gen",
        "duration": null,
        "refund_credits": null,
        "stream": true,
        "error_type": null,
        "error_message": null
      },
      "is_liked": false,
      "user_id": "284ac3ca-a2cf-4b0c-a1b5-64d7315fbe28",
      "is_trashed": false,
      "reaction": null,
      "created_at": "2024-03-26T07:01:44.235Z",
      "status": "submitted",
      "title": "",
      "play_count": 0,
      "upvote_count": 0,
      "is_public": false
    }
  ],
  "metadata": {
    "tags": null,
    "prompt": "",
    "gpt_description_prompt": "an atmospheric metal song about dancing all night long",
    "audio_prompt_id": null,
    "history": null,
    "concat_history": null,
    "type": "gen",
    "duration": null,
    "refund_credits": null,
    "stream": true,
    "error_type": null,
    "error_message": null
  },
  "major_model_version": "v3",
  "status": "running",
  "created_at": "2024-03-26T07:01:44.214Z",
  "batch_size": 2
}
```
#### 4. 注意事项
- 本项目仅用于学习交流，不得用于商业用途
- 个人账号部分功能无法使用