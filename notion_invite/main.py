from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

# ブラウザを起動
driver = webdriver.Chrome(executable_path='/opt/homebrew/bin/chromedriver')

# Notionの該当ページにアクセス
driver.get('******')

# Shareボタンをクリック（このセレクタは仮のものなので、実際のページに合わせて変更してください）
share_button = WebDriverWait(driver, 10).until(
    EC.presence_of_element_located((By.CSS_SELECTOR, ".share-button-css-selector"))
)
share_button.click()

# メールアドレス入力フィールドにメールアドレスを入力
email_input = WebDriverWait(driver, 10).until(
    EC.presence_of_element_located((By.CSS_SELECTOR, ".email-input-css-selector"))
)
email_input.send_keys('******')

# Invite ボタンをクリック
invite_button = WebDriverWait(driver, 10).until(
    EC.presence_of_element_located((By.CSS_SELECTOR, ".invite-button-css-selector"))
)
invite_button.click()

# 作業終了後、ブラウザを閉じる
driver.quit()
