import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  try {
    await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page.waitForTimeout(2000);

    // 检查搜索按钮是否存在
    const button = await page.$('button.btn-primary');
    if (button) {
      const isVisible = await button.isVisible();
      const text = await button.textContent();
      const boundingBox = await button.boundingBox();
      console.log('✅ 搜索按钮存在');
      console.log('   文本:', text?.trim());
      console.log('   可见:', isVisible);
      console.log('   位置:', boundingBox);
    } else {
      console.log('❌ 搜索按钮不存在');
    }

    // 检查 input-group 内容
    const inputGroup = await page.$('.input-group');
    if (inputGroup) {
      const children = await inputGroup.$$('*');
      console.log('\n.input-group 包含', children.length, '个子元素:');
      for (let i = 0; i < children.length; i++) {
        const tag = await page.evaluate(el => el.tagName, children[i]);
        const cls = await page.evaluate(el => el.className, children[i]);
        console.log(`  [${i}] <${tag}> class="${cls}"`);
      }
    }
  } finally {
    await browser.close();
  }
})();
