#!/bin/bash
set -e

echo "Installing markitdown..."
pip install --no-cache-dir markitdown[all]

echo "Starting markitdown server..."
# markitdown 提供了一个内置的 HTTP 服务器
# 如果没有内置服务器，我们用 Flask 包装一个简单的 API
pip install --no-cache-dir flask

cat > /app.py << 'EOF'
from flask import Flask, request, jsonify
from markitdown import MarkItDown
import tempfile
import os

app = Flask(__name__)
md = MarkItDown()

@app.route('/health', methods=['GET'])
def health():
    return jsonify({"status": "ok"})

@app.route('/convert', methods=['POST'])
def convert():
    if 'file' not in request.files:
        return jsonify({"error": "No file provided"}), 400

    file = request.files['file']
    if file.filename == '':
        return jsonify({"error": "No file selected"}), 400

    # 保存临时文件
    suffix = os.path.splitext(file.filename)[1]
    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
        file.save(tmp.name)
        tmp_path = tmp.name

    try:
        result = md.convert(tmp_path)
        return jsonify({
            "text": result.text_content,
            "filename": file.filename
        })
    except Exception as e:
        return jsonify({"error": str(e)}), 500
    finally:
        os.unlink(tmp_path)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081)
EOF

exec python /app.py
