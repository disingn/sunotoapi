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
创建音乐任务
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
支持如下参数：
- 默认参数
    - gpt_description_prompt: 生成音乐的描述
    - mv: 音乐模型
    - prompt: 生成音乐的提示
    - make_instrumental: 是否生成无人声音乐
```json
{
    "gpt_description_prompt": "an atmospheric metal song about dancing all night long",
    "mv": "chirp-v3-0",
    "prompt": "",
    "make_instrumental": false
}
```
- 自定义
    - prompt: 歌词
    - tags: 音乐标签
    - mv: 音乐模型
    - title: 音乐标题
    - continue_clip_id: 继续生成音乐的clip_id
    - continue_at: 继续生成音乐的时间
    
```json
{
  "prompt": "[Verse]\nEvery morning, when I wake up\nI stumble to the kitchen to get my cup (cup)\nThe smell, the taste, it's like a dream\nI'm addicted to that caffeinated beam\n\n[Chorus]\nI need my java fix, it's my daily high (high)\nGotta have my coffee, don't ask me why (why)\nBrew it strong, brew it black, can't get enough\nThat sweet, dark liquid, it keeps me buzzin' (buzzin')\n\n[Verse 2]\nEspresso, latte, cappuccino too\nI'll take it any way, as long as it's brew\nFrom the fancy cafes to the corner shops\nI'm on a mission to find the perfect crop (yeah)",
  "tags": "epic blues",
  "mv": "chirp-v3-0",
  "title": "Coffee Addiction",
  "continue_clip_id": null,
  "continue_at": null
}
```
- 自定义纯音乐
    - prompt: 歌词
    - tags: 音乐标签
    - mv: 音乐模型
    - title: 音乐标题
    - continue_clip_id: 继续生成音乐的clip_id
    - continue_at: 继续生成音乐的时间
```json
{
  "prompt": "",
  "tags": "epic blues",
  "mv": "chirp-v3-0",
  "title": "Coffee Addiction",
  "continue_clip_id": null,
  "continue_at": null
}
```
查询音乐任务
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
使用ai生成歌词
```shell
curl --location --request POST 'localhost:3560/v2/lyrics/create' \
--header 'Content-Type: application/json' \
--data-raw '{"prompt":""}'
```
响应json：
```json
{
  "id": "4ad435dd-b3f1-4ed7-b316-1868ff4ffe55"
}
```
查询生成的歌词
```shell
curl --location --request POST 'localhost:3560/v2/lyrics/task' \
--header 'Content-Type: application/json' \
--data-raw '{
    "ids":"4ad435dd-b3f1-4ed7-b316-1868ff4ffe55"
}'
```
响应json：
```json
{
  "text": "[Verse]\nI saw you sippin' on your latte with grace\nAcross the room, I couldn't help but gaze (gazin')\nYour smile was like a sunbeam on a cloudy day\nIn that moment, I knew I had to find a way\n\n[Chorus]\nCoffee shop love affair, brewing in the air\nAin't nothin' like the feeling when we're there (when we're there)\nCoffee shop love affair, hearts are gonna dare\nTo catch a glimpse of something rare\n\n[Verse 2]\nI ordered my Americano, extra light\nAnd slowly made my way closer in sight (closer in sight)\nOur eyes met, and the world began to fade\nIn that coffee shop, our love was made (made, baby)",
  "title": "The Coffee Shop Love Affair",
  "status": "complete"
}
```
兼容 openai /v1/chat/completions格式
```shell
curl --location --request POST 'localhost:3560/v1/chat/completions' \
--header 'Content-Type: application/json' \
--data-raw '{
    "model": "chirp-v3-0",
    "messages": [
        {
            "role": "user",
            "content": "制作歌曲《万能的mj）中文歌词 mj是-个群体的统称"
        }
    ]
}'
```
响应json：
```json
{
  "choices": [
    {
      "finish_reason": "stop",
      "index": 0,
      "logprobs": null,
      "message": {
        "content": "xxxx",
        "role": "assistant"
      }
    }
  ],
  "created": 1711612671,
  "id": "chatcmpl-7QyqpwdfhqwajicIEznoc6Q47XAyW",
  "model": "chirp-v3-0",
  "object": "chat.completion",
  "usage": {
    "completion_tokens": 17,
    "prompt_tokens": 57,
    "total_tokens": 74
  }
}
```
#### 4. 注意事项
- 本项目仅用于学习交流，不得用于商业用途
- 个人账号部分功能无法使用