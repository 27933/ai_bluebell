import { chromium } from '@playwright/test';

(async () => {
  console.log('启动浏览器，开始截图...\n');

  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();

  try {
    // 设置视口大小
    await page.setViewportSize({ width: 1200, height: 800 });

    // 访问首页
    console.log('正在访问 http://localhost:5173...');
    await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });

    // 等待页面加载完成
    await page.waitForTimeout(2000);

    // 截图
    const screenshotPath = '/tmp/home-page.png';
    await page.screenshot({ path: screenshotPath, fullPage: false });
    console.log(`✅ 首页截图已保存: ${screenshotPath}`);

    // 获取页面信息
    const title = await page.title();
    console.log(`\n📄 页面标题: ${title}`);

    // 检查页面元素
    const navbarExists = (await page.$('.navbar')) !== null;
    const articleCards = await page.$$('.article-card');
    const tags = await page.$$('.tag');

    console.log('\n📋 页面元素检查:');
    console.log(`  ✅ 导航栏: ${navbarExists ? '✓ 存在' : '✗ 不存在'}`);
    console.log(`  ✅ 文章卡片: ${articleCards.length > 0 ? `✓ ${articleCards.length} 个` : '✗ 不存在'}`);
    console.log(`  ✅ 标签: ${tags.length > 0 ? `✓ ${tags.length} 个` : '✗ 不存在'}`);

    // 获取导航栏信息
    const brandText = await page.textContent('.navbar-brand');
    console.log(`\n🏷️  品牌文本: "${brandText?.trim()}"`);

    // 获取搜索框
    const searchInput = await page.$('input[placeholder*="搜索"]');
    if (searchInput) {
      const placeholder = await searchInput.getAttribute('placeholder');
      console.log(`\n🔍 搜索框: "${placeholder}"`);
    }

    // 获取计算后的样式信息
    const navbarStyle = await page.evaluate(() => {
      const navbar = document.querySelector('.navbar');
      if (!navbar) return null;
      const computed = window.getComputedStyle(navbar);
      return {
        backgroundColor: computed.backgroundColor,
        padding: computed.padding,
        boxShadow: computed.boxShadow,
        position: computed.position,
      };
    });

    if (navbarStyle) {
      console.log('\n🎨 Navbar样式:');
      console.log(`  背景色: ${navbarStyle.backgroundColor}`);
      console.log(`  Padding: ${navbarStyle.padding}`);
      console.log(`  阴影: ${navbarStyle.boxShadow}`);
      console.log(`  定位: ${navbarStyle.position}`);
    }

    // 获取文章卡片样式
    const cardStyle = await page.evaluate(() => {
      const card = document.querySelector('.article-card');
      if (!card) return null;
      const computed = window.getComputedStyle(card);
      return {
        borderRadius: computed.borderRadius,
        boxShadow: computed.boxShadow,
        backgroundColor: computed.backgroundColor,
        transition: computed.transition,
      };
    });

    if (cardStyle) {
      console.log('\n🎨 文章卡片样式:');
      console.log(`  圆角: ${cardStyle.borderRadius}`);
      console.log(`  阴影: ${cardStyle.boxShadow}`);
      console.log(`  背景色: ${cardStyle.backgroundColor}`);
      console.log(`  过渡: ${cardStyle.transition}`);
    }

    // 获取标签样式
    const tagStyle = await page.evaluate(() => {
      const tag = document.querySelector('.tag');
      if (!tag) return null;
      const computed = window.getComputedStyle(tag);
      return {
        backgroundColor: computed.backgroundColor,
        color: computed.color,
        padding: computed.padding,
        borderRadius: computed.borderRadius,
      };
    });

    if (tagStyle) {
      console.log('\n🎨 标签样式:');
      console.log(`  背景色: ${tagStyle.backgroundColor}`);
      console.log(`  文字色: ${tagStyle.color}`);
      console.log(`  Padding: ${tagStyle.padding}`);
      console.log(`  圆角: ${tagStyle.borderRadius}`);
    }

    console.log('\n✅ 页面分析完成！');
    console.log(`📸 截图已保存到: ${screenshotPath}\n`);

  } catch (error) {
    console.error('❌ 出错:', error.message);
  } finally {
    await browser.close();
  }
})();
