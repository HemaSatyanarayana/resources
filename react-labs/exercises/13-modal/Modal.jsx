import { useEffect } from 'react'
import { createPortal } from 'react-dom'

/**
 * Modal — render outside the parent DOM tree with a portal, close on Escape /
 * backdrop click / close button.
 *
 * Props:
 *   isOpen    boolean
 *   onClose   () => void
 *   children  modal content
 *
 * Required UI (the tests rely on this):
 *   - renders nothing when isOpen is false
 *   - when open: a backdrop element with data-testid="backdrop", containing an
 *     element with role="dialog" (aria-modal="true") that wraps `children` and
 *     a button labelled "Close"
 *   - the whole thing is portaled into document.body
 *
 * Rules:
 *   - Pressing Escape calls onClose (listener on document; added only while open,
 *     removed on cleanup).
 *   - Clicking the backdrop itself (not its content) calls onClose.
 *   - Clicking the Close button calls onClose.
 */
export default function Modal({ isOpen, onClose, children }) {
  useEffect(() => {
    if (!isOpen) return
    // TODO: add a keydown listener to document that calls onClose on Escape.
    //       Return a cleanup that removes it.
  }, [isOpen, onClose])

  if (!isOpen) return null

  const onBackdropClick = (e) => {
    // TODO: only close when the click target IS the backdrop itself
    //       (e.target === e.currentTarget), so clicks inside the dialog don't close it.
  }

  return createPortal(
    <div data-testid="backdrop" onClick={onBackdropClick}>
      <div role="dialog" aria-modal="true">
        {children}
        <button onClick={onClose}>Close</button>
      </div>
    </div>,
    document.body
  )
}
