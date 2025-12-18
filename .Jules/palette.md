## 2024-05-22 - [Format Selection & Toggle Logic]
**Learning:** Changing interactive `div`s to `button`s is the single most impactful accessibility win for lists of choices. It grants keyboard focus and activation for free, but requires careful styling (resetting user agent styles) to match the original "card" look.
**Action:** Always prefer `<button>` or `<input type="radio">` for selection lists. Use `w-100 text-start` (Bootstrap) to make buttons behave like block-level list items.

**Learning:** When toggling visibility with JS, always check explicitly for the `none` state rather than relying on falsy checks of `style.display`, especially if the element is visible by default (where `style.display` is empty string).
**Action:** Use `if (el.style.display === 'none')` for show logic, and `else` for hide logic, to handle initial visible states correctly.
