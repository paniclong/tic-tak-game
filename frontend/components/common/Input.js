import React from "react";

const Input = ({ type, value, isChecked, text, onChange }) => {
  return (
    <div>
      <span>
      {text}
      </span>
      <input type={type} value={value} checked={isChecked} onChange={onChange}>
      </input>
    </div>
  );
}

export default Input
