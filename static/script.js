let ws;
let selectedEmotion = '';

// 화면 전환 함수
function showScreen(screenId) {
    document.querySelectorAll('.screen').forEach(screen => {
        screen.classList.remove('active');
    });
    document.getElementById(screenId).classList.add('active');
}

// WebSocket 연결 함수
function connectWebSocket() {
    ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {
        console.log('WebSocket 연결됨');
        ws.send(selectedEmotion);
    };

    ws.onmessage = (event) => {
        const message = event.data;
        console.log(message)
        if (message === 'matched') {
            // When a match is found, wait 2 seconds before showing the chat screen
            setTimeout(() => {
                // 채팅 메시지 초기화
                const messagesDiv = document.getElementById('chat-messages');
                messagesDiv.innerHTML = ''; // 모든 메시지 제거
                showScreen('chat-screen');
            }, 2000);
        } else {
            addMessage(message, 'received');
        }
    };


    ws.onclose = () => {
        console.log('WebSocket 연결 끊김');
        // 채팅 기록 초기화
        const messagesDiv = document.getElementById('chat-messages');
        messagesDiv.innerHTML = '';
        showScreen('emotion-select');
    };

}

// 메시지 추가 함수
function addMessage(message, type) {
    const messagesDiv = document.getElementById('chat-messages');
    const messageElement = document.createElement('div');
    messageElement.className = `message ${type}`;
    messageElement.textContent = message;
    messagesDiv.appendChild(messageElement);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

// 이벤트 리스너 설정
document.addEventListener('DOMContentLoaded', () => {
    // 감정 버튼 클릭 이벤트
    document.querySelectorAll('.emotion-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            selectedEmotion = btn.dataset.emotion;
            document.getElementById('selected-emotion-text').textContent =
                selectedEmotion === 'happy' ? '기쁨 😊' :
                    selectedEmotion === 'angry' ? '화남 😠' : '슬픔 😢';
            showScreen('matching-screen');
            connectWebSocket();
        });
    });

    // 매칭 취소 버튼
    document.getElementById('cancel-matching').addEventListener('click', () => {
        if (ws) {
            ws.close();
        }
        showScreen('emotion-select');
    });

    // 채팅방 나가기 버튼
    document.getElementById('leave-chat').addEventListener('click', () => {
        if (ws) {
            ws.close();
        }
        showScreen('emotion-select');
    });

    // 메시지 전송
    const sendMessage = () => {
        const input = document.getElementById('message-input');
        const message = input.value.trim();
        if (message && ws) {
            ws.send(message);
            addMessage(message, 'sent');
            input.value = '';
        }
    };

    // 전송 버튼 클릭
    document.getElementById('send-message').addEventListener('click', sendMessage);

    // Enter 키 입력
    document.getElementById('message-input').addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            sendMessage();
        }
    });
});