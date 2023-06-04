import { Component } from 'react';

function Media({ src, alt }) {
  if (src)
    return <img src={src} alt={alt} />
  else
    return <div />
}

class Story extends Component {
  constructor(props) {
    super(props)
    this.state = {
      title: "",
      excerpt: "",
      storyLink: "",
      media: null,
      caption: null,
    }
    this.alreadyFetched = false
  }
  async componentDidMount() {
    if (this.alreadyFetched)
      return // see https://stackoverflow.com/questions/71755119/reactjs-componentdidmount-executes-twice
    this.alreadyFetched = true
    const response = await fetch('https://stories.workingclasshistory.com/api/v1/one_random_from_today')
    // you may want to add some better error handling here.
    if (response.ok) {
      const data = await response.json()
      const { title, excerpt, url, media } = data
      this.setState({ title, excerpt, storyLink: url })
      if (media) {
        const { url, caption } = media
        this.setState({ media: url, caption })
      }
    } else {
      throw new Error("error fetching data")
    }
  }
  render() {
    const { title, excerpt, media, caption, storyLink } = this.state
    return (
      <div className="App">
        <h2>{title}</h2>
        <Media src={media} alt={caption} />
        <div>{excerpt}</div>
        <a href={storyLink} target="_blank" rel="noreferrer">Read more...</a>
      </div>
    )
  }
}

function App() {
  return <Story />
}

export default App;
