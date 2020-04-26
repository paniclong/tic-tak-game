import React from "react";

const Button = ({ value, onClick, className, styles, disabled }) => {
  return (
    <button className={className} onClick={onClick} style={styles} disabled={disabled}>
      {value}
    </button>
  );
}

export default Button
