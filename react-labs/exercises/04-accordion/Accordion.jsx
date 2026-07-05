import { useState } from "react";

/**
 * Accordion — conditional rendering + single vs multi-open state.
 *
 * Props:
 *   items          [{ id, title, content }]
 *   allowMultiple  boolean  if false (default), opening one closes the others
 *
 * Required UI (the tests rely on this):
 *   - one header <button> per item, whose text is the item's title
 *   - each header button has aria-expanded="true"|"false"
 *   - an item's `content` is only in the document while that item is open
 *
 * Rules:
 *   - Clicking a header toggles that item.
 *   - allowMultiple=false: at most one item open at a time.
 *   - allowMultiple=true: items open/close independently.
 */
export default function Accordion({ items = [], allowMultiple = false }) {
  const [openIds, setOpenIds] = useState([]);

  const toggle = (id) => {
    // TODO: update openIds.
    //  - if the id is already open, close it.
    //  - otherwise open it: append when allowMultiple, else replace with [id].
    setOpenIds((prevIds) => {
      if (prevIds.includes(id)) {
        return prevIds.filter((openId) => openId !== id);
      } else {
        return allowMultiple ? [...prevIds, id] : [id];
      }
    });
  };

  return (
    <div>
      {items.map((item) => {
        const isOpen = openIds.includes(item.id);
        return (
          <div key={item.id}>
            <button aria-expanded={isOpen} onClick={() => toggle(item.id)}>
              {item.title}
            </button>
            {isOpen && <div role="region">{item.content}</div>}
          </div>
        );
      })}
    </div>
  );
}
