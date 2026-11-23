import { describe, expect, test } from 'bun:test'
import { uniqueCategories } from '../app/utils'

describe('policy source categories', () => {
  test('merges the categories of all sources once and sorted', () => {
    expect(uniqueCategories([
      { categories: ['Pod Security', 'Best Practices'] },
      { categories: ['Best Practices', 'Multi-Tenancy'] },
    ])).toEqual(['Best Practices', 'Multi-Tenancy', 'Pod Security'])
  })

  test('handles missing sources and categories', () => {
    expect(uniqueCategories(undefined)).toEqual([])
    expect(uniqueCategories(null)).toEqual([])
    expect(uniqueCategories([{ categories: null }, {}])).toEqual([])
  })
})
