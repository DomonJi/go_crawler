import React, { Component } from 'react'
import TextField from 'material-ui/TextField'
import Button from 'material-ui/Button'
import List, { ListItem } from 'material-ui/List'
import Card, { CardContent } from 'material-ui/Card'
import './App.css'

class App extends Component {
  constructor(props){
    super(props)
    this.state = {
      search: '',
      result: [],
    }
  }

  handleChange = event => {
    this.setState({ search: event.target.value })
  }

  handleKeyPress = event => {
    event.key === 'Enter' && this.search()
  }

  search = () => {
    if (!this.state.search) return
    const searchValue = this.state.search
    fetch(`http://localhost:9200/spider/_search?q=${encodeURIComponent(searchValue)}`, {
      method: 'GET',
    })
    .then(res => res.json())
    .then(res => {
      if (res && res.hits) this.setState({
        result: res.hits.hits
      })
    })
    .catch(console.log.bind(console))
  }

  render() {
    return (
      <div className="App">
        <TextField
          id="search"
          name="search"
          value={this.state.search}
          onChange={this.handleChange}
          onKeyPress={this.handleKeyPress}
          autoFocus
        />
        <Button id="search-button" onClick={this.search}>Search</Button>
        <List>
          {this.state.result.map(item => (
            <ListItem>
              <Card className="result-card" onClick={() => window.open(item._source.url)}>
                <CardContent>
                  <h2>{item._source.name}</h2>
                  <p>{item._source.summary}</p>
                </CardContent>
              </Card>
            </ListItem>
          ))}
        </List>
      </div>
    )
  }
}

export default App
