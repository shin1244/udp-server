# UDP Pixel Streamer 📡

Go와 Ebiten 게임 엔진을 활용하여 UDP로 전송된 픽셀 데이터를 실시간으로 시각화하는 연습용 프로젝트입니다.

## 📝 프로젝트 소개 (Project Overview)

이 프로젝트는 네트워크 프로그래밍과 Go의 동시성 패턴(Goroutine, Channel)을 학습하기 위해 만들어졌습니다.
외부에서 UDP 패킷으로 픽셀 정보(좌표, 색상)를 보내면, 뷰어에서 이를 수신하여 실시간으로 이미지를 재구성합니다.

현재는 기본적인 수신 및 렌더링 기능이 구현되어 있으며, 추후 신뢰성 있는 UDP 통신(Reliable UDP)을 위한 재전송 로직을 추가할 예정입니다.

## 🛠 사용 기술 (Tech Stack)

- **Language**: Go (Golang)
- **Library**: [Ebiten v2](https://github.com/hajimehoshi/ebiten) (2D Game Engine)
- **Network**: UDP Socket
- **Concurrency**: Goroutines, Channels

## 🚀 기능 (Features)

- **UDP Server**: `:8080` 포트에서 UDP 패킷 수신
- **Concurrency Pipeline**:
  - `Receiver`: 패킷 수신 및 파싱
  - `Packet Channel`: 데이터 버퍼링 및 전달
  - `Updater`: 픽셀 데이터 갱신 (Thread-Safe)
- **Rendering**: Ebiten을 통한 256x256 픽셀 실시간 드로잉

## 🔮 향후 계획 (To-Do)

- [O] **패킷 유실 복구 (Packet Recovery)**: 순서가 바뀌거나 유실된 패킷을 감지하고 재전송을 요청하는 로직 구현 (ARQ)
- [ ] **전송 최적화**: 픽셀 단위 전송 방식에서 청크(Chunk) 단위 전송으로 변경하여 오버헤드 감소화
