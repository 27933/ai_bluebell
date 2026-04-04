import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  try {
    await page.setViewportSize({ width: 1600, height: 900 });
    await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page.waitForTimeout(2000);

    // 检查按钮
    const button = await page.$('.btn-primary');
    if (button) {
      const text = await button.textContent();
      const boundingBox = await button.boundingBox();
      const display = await button.evaluate(el => window.getComputedStyle(el).display);
      const visibility = await button.evaluate(el => window.getComputedStyle(el).visibility);
      console.log('✅ 按钮存在');
      console.log('   文本:', text?.trim());
      console.log('   display:', display);
      console.log('   visibility:', visibility);
      console.log('   位置:', boundingBox);
    } else {
      console.log('❌ 按钮不存在于 DOM');
    }

    // 检查 input-group
    const inputGroup = await page.$('.input-group');
    if (inputGroup) {
      const html = await inputGroup.innerHTML();
      console.log('\ninput-group HTML:');
      console.log(html.substring(0, 200));
    }
  } finally {
    await browser.close();
  }
})();
