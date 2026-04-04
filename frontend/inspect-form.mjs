import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  try {
    await page.setViewportSize({ width: 1600, height: 900 });
    await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page.waitForTimeout(2000);

    // 检查搜索表单各元素
    const col = await page.$('.col-md-8');
    const inputGroup = await page.$('.input-group');
    const input = await page.$('.form-control');
    const button = await page.$('.btn-primary');

    console.log('col-md-8:', await col?.boundingBox());
    console.log('input-group:', await inputGroup?.boundingBox());
    console.log('input:', await input?.boundingBox());
    console.log('button:', await button?.boundingBox());

  } finally {
    await browser.close();
  }
})();
