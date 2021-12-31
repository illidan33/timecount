import os
import time
from http.server import HTTPServer, BaseHTTPRequestHandler

import cv2
import face_recognition
import numpy as np


def capture():
    try:
        video_capture = cv2.VideoCapture(0)

        obama_image = face_recognition.load_image_file(source_pic)
        obama_face_encodings = face_recognition.face_encodings(obama_image)
        if len(obama_face_encodings) == 0:
            print("图片无法有效识别人脸")
            return False

        # 取第一个人脸信息
        known_face_encodings = [
            obama_face_encodings[0],
        ]
        known_face_names = [
            your_name,
        ]

        start = time.time() * 1000
        while True:
            now = time.time() * 1000
            if now - start >= 1000:
                break
            ret, frame = video_capture.read()
            small_frame = cv2.resize(frame, (0, 0), fx=0.25, fy=0.25)
            rgb_small_frame = small_frame[:, :, ::-1]
            face_locations = face_recognition.face_locations(rgb_small_frame, 3)
            face_encodings = face_recognition.face_encodings(rgb_small_frame, face_locations)

            face_names = []
            for face_encoding in face_encodings:
                matches = face_recognition.compare_faces(known_face_encodings, face_encoding)
                name = "Unknown"
                face_distances = face_recognition.face_distance(known_face_encodings, face_encoding)
                best_match_index = np.argmin(face_distances)
                if matches[best_match_index]:
                    name = known_face_names[best_match_index]

                face_names.append(name)

            if your_name in face_names:
                return True
        return False
    except Exception:
        print("internal error")
    finally:
        video_capture.release()
        cv2.destroyAllWindows()


class Request(BaseHTTPRequestHandler):
    def do_GET(self):
        rs = False
        if self.path == "/":
            rs = capture()
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        if rs:
            self.wfile.write('ok'.encode())
        else:
            self.wfile.write('false'.encode())


if __name__ == "__main__":
    source_pic = "source.jpg"  # 替换你的图片地址
    your_name = "huoenhaimu"  # 可不修改
    if not os.path.exists(source_pic):
        print(source_pic, " not found")
        exit()
    try:
        server = HTTPServer(("127.0.0.1", 7001), Request)
        server.serve_forever()
    except KeyboardInterrupt:
        print("exit")
