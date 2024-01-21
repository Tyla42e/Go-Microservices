import React, { useState } from "react";
//import axios from "axios";

const PostCreate = () => {
  const [title, setTitle] = useState("");

  const onSubmit = (event) => {
    event.preventDefault();
    console.log("Title: " + title)
    console.log("Calling create")
    // await axios.post("http://mywibbleposts.com/create", {
    //   title,
    // });

    fetch('http://mywibbleposts.com/posts/create', {
      method: 'post',
      mode: 'cors',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify({"title": title}) 
     });
    setTitle("");
  };

  return (
    <div>
      <form onSubmit={onSubmit}>
        <div className="form-group">
          <label>Title</label>
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="form-control"
          />
        </div>
        <button className="btn btn-primary">Submit</button>
      </form>
    </div>
  );
};

export default PostCreate;
