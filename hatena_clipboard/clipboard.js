<script>
  document.addEventListener('DOMContentLoaded', function() {
    // コードブロックを取得
    var codeBlocks = document.querySelectorAll('pre.code');
    
    // すべてのコードブロックにコピーボタンを追加
    for (var i = 0; i < codeBlocks.length; i++) {
      var copyButton = document.createElement('button');
      copyButton.className = 'copy-button';
      copyButton.textContent = 'Copy code';
      copyButton.onclick = function() {
        var codeElem = this.parentNode.querySelector('code') || this.parentNode;
        var textArea = document.createElement('textarea');
        textArea.value = codeElem.textContent.replace(/Copy code$/, ''); // "Copy code" テキストを削除
        document.body.appendChild(textArea);
        textArea.select();
        document.execCommand('copy');
        document.body.removeChild(textArea);
      }
      codeBlocks[i].appendChild(copyButton);
    }
  });
</script>
