import { describe, expect, test } from 'bun:test'
import { renderPolicyDescription } from '../app/utils/markdown'

describe('policy descriptions', () => {
  test('renders documentation formatting', () => {
    const html = renderPolicyDescription([
      '## Configuration',
      '',
      '- Set `replicas` to **3**',
      '  - Enable *TLS*',
      '',
      '1. Read [the docs](https://example.com/docs)',
      '',
      '> Required configuration',
      '',
      '```yaml',
      'replicas: 3',
      '```',
    ].join('\n'))

    expect(html).toContain('<h2>Configuration</h2>')
    expect(html).toContain('<code>replicas</code>')
    expect(html).toContain('<strong>3</strong>')
    expect(html).toContain('<em>TLS</em>')
    expect(html.match(/<ul>/g)).toHaveLength(2)
    expect(html).toContain('<ol>')
    expect(html).toContain('<a href="https://example.com/docs">the docs</a>')
    expect(html).toContain('<blockquote>')
    expect(html).toContain('<pre><code class="language-yaml">replicas: 3\n</code></pre>')
  })

  test('preserves plain text, Unicode, and single line breaks', () => {
    expect(renderPolicyDescription('')).toBe('')
    expect(renderPolicyDescription('Enable TLS')).toBe('<p>Enable TLS</p>\n')
    expect(renderPolicyDescription('設定 ✓\nEnable TLS')).toBe('<p>設定 ✓<br>\nEnable TLS</p>\n')
    expect(renderPolicyDescription('A & B')).toBe('<p>A &amp; B</p>\n')
    expect(renderPolicyDescription('https://example.com')).not.toContain('<a')
  })

  test('keeps HTML literal inside code', () => {
    const html = renderPolicyDescription('`<script>alert(1)</script>`\n\n```html\n<img src=x onerror=alert(1)>\n```')
    expect(html).toContain('<code>&lt;script&gt;alert(1)&lt;/script&gt;</code>')
    expect(html).toContain('&lt;img src=x onerror=alert(1)&gt;')
    expect(html).not.toContain('<script')
    expect(html).not.toContain('<img')
  })

  test.each([
    '<script>alert(1)</script>',
    '<img src=x onerror=alert(1)>',
    '<svg onload=alert(1)></svg>',
    '<iframe src="javascript:alert(1)"></iframe>',
    '<a href="https://example.com" onclick="alert(1)">link</a>',
  ])('escapes raw HTML: %s', (source) => {
    const html = renderPolicyDescription(source)
    expect(html).toContain('&lt;')
    expect(html).not.toMatch(/<(?:script|img|svg|iframe|a)\b/i)
  })

  test.each([
    'javascript:alert(1)',
    'JaVaScRiPt:alert(1)',
    'jav&#x61;script:alert(1)',
    'javascript&#58;alert(1)',
    'vbscript:msgbox(1)',
    'data:text/html;base64,PHNjcmlwdD4=',
  ])('rejects unsafe link and image destinations: %s', (url) => {
    expect(renderPolicyDescription(`[link](${url})`)).not.toContain('<a ')
    expect(renderPolicyDescription(`![image](${url})`)).not.toContain('<img ')
  })
})
