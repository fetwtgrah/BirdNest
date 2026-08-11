package utils

func EmailContext() string {
	return `
<div style="max-width: 440px; margin: 40px auto; padding: 0; font-family: -apple-system, 'Segoe UI', 'Microsoft YaHei', sans-serif; background-color: #FFF8F0;">
    <div style="text-align: center; padding: 36px 20px 20px;">
        <div style="font-size: 36px;">🐦</div>
        <p style="font-size: 13px; letter-spacing: 3px; color: #C88A5A; margin: 8px 0 0; font-weight: 600;">BIRDNEST</p>
    </div>
    <div style="padding: 0 32px 40px;">
        <p style="font-size: 15px; color: #6B4226; text-align: center; margin: 0 0 24px;">Hi %s，欢迎回巢 🌿</p>
        <div style="border: 2px dashed #E0B482; border-radius: 24px; padding: 28px 20px; text-align: center; background: #FFFDF9;">
            <p style="font-size: 12px; color: #C88A5A; margin: 0 0 12px; letter-spacing: 1px;">你的验证码</p>
            <p style="font-size: 38px; font-weight: 700; letter-spacing: 8px; color: #8B5E3C; margin: 0; font-family: 'Courier New', monospace;">%s</p>
        </div>
        <p style="font-size: 13px; color: #B89073; text-align: center; line-height: 1.8; margin-top: 24px;">
            验证码 5 分钟内有效，用完请让它安心过期 🌙<br>
            请勿告诉任何人，巢里的秘密只留给你自己
        </p>
    </div>
    <div style="text-align: center; padding: 16px; border-top: 1px solid #F0DFC8;">
        <p style="font-size: 11px; color: #D4B896; margin: 0;">BirdNest · 一个温暖的小窝</p>
    </div>
</div>
`
}
