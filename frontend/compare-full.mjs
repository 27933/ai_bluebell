import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  try {
    // 设置更大的视口
    await page.setViewportSize({ width: 1280, height: 900 });
    
    console.log('🔄 正在加载当前前端（完整视图）...');
    await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page.waitForTimeout(2000);
    await page.screenshot({ path: '/tmp/full-current.png', fullPage: true });
    console.log('✅ 完整前端截图已保存: /tmp/full-current.png');

  } catch (error) {
    console.error('❌ 出错:', error.message);
  } finally {
    await browser.close();
  }
})();
