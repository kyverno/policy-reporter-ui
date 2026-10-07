import MarkdownIt from 'markdown-it'

// Policy descriptions are untrusted. Keep raw HTML disabled and retain the
// parser's URL validation; only its escaped output may be rendered as HTML.
const markdown = new MarkdownIt({
  html: false,
  breaks: true,
  linkify: false,
  typographer: false,
})

export const renderPolicyDescription = (description: string): string => markdown.render(description)
