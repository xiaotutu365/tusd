## TUSD新增接口说明

### 基础介绍

#### 启动命令
```shell
tusd -s3-bucket=my-test-bucket.com 
    -s3-endpoint https://mystoreage.example.com
    -hooks-enabled-events pre-download,pre-create
    -hooks-http http://localhost:8080/write
    -hooks-http-forward-headers Upload-Metadata
```
* `-s3-bucket`为S3的bucket，`s3-endpoint`为S3的端点；
* `-hooks-enabled-events`代表激活的钩子，这里表示激活的是pre-download,pre-create钩子
* `-hooks-http`表示钩子实际的请求URL
* `-hooks-http-forward-headers`表示上传或者下载过程中，元数据传递的参数名称

#### S3服务启动配置
```shell
tusd -s3-bucket=my-test-bucket.com -s3-endpoint https://mystoreage.example.com
```

#### Disk服务启动配置
```shell
tusd -upload-dir=./uploads
```

### 新增操作接口

#### 上传客户端
```java
public class Client {

    public static void main(String[] args) throws ProtocolException, IOException {
        final TusClient client = new TusClient();
        client.setUploadCreationURL(new URL("http://localhost:8080/files/"));
        client.enableResuming(new TusURLMemoryStore());
        File file = new File("/path/file");
        final TusUpload upload;
        try {
            upload = new TusUpload(file);
            Map<String, String> metadata = new HashMap<>();
            // 元数据
            metadata.put("token", "ab63d270-7edc-4338-8b4d-01866e3c37a1");
            metadata.put("project", "catl_order");
            upload.setMetadata(metadata);
        } catch (FileNotFoundException e) {
            throw new RuntimeException(e);
        }

        System.out.println("Starting upload...");

        TusExecutor executor = new TusExecutor() {
            @Override
            protected void makeAttempt() throws ProtocolException, IOException {
                TusUploader uploader = client.resumeOrCreateUpload(upload);
                uploader.setChunkSize(1024);
                do {
                    long totalBytes = upload.getSize();
                    long bytesUploaded = uploader.getOffset();
                    double progress = (double) bytesUploaded / totalBytes * 100;

                    System.out.printf("Upload at %06.2f%%.\n", progress);
                } while (uploader.uploadChunk() > -1);
                uploader.finish();
                System.out.println("Upload finished.");
                System.out.format("Upload available at: %s", uploader.getUploadURL().toString());
            }
        };
        executor.makeAttempts();
    }
}
```

#### 删除接口
URL: http://[::]:8080/files/{uuid}
Method: DELETE
其中{uuid}为上传返回的uuid值。


### 批量接口下载
URL: http://[::]:8080/download
Method: POST
HEADER:
    Download-Metadata: base64(filename: {uuid1},{uuid2},...,{uuid3})


### 下载钩子
URL: http://[::]:8080/files/{uuid}
Method: Get

Cli: pre-download: {hook-url}
其中hook-url代表前置钩子


### 多存储类型
URL: http://[::]:8080/files
Method: Post

HEADER: upload-type: s3,disk
其中s3表示存储在aws s3中，disk表示存储在本地磁盘，该数据是存储在`Upload-Metadata`中的，以KV值来存储。