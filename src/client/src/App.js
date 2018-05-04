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
    this.handleChange = this.handleChange.bind(this)
    this.search = this.search.bind(this)
  }

  handleChange(event) {
    this.setState({ search: event.target.value })
  }

  search(){
    if (!this.state.search) return
    const searchValue = this.state.search
    fetch(`http://localhost:9200/spider/_search?q=${encodeURIComponent(searchValue)}`, {
      method: 'GET',
    }).then(res => res.json())
    .then(res => this.setState({ result: res.hits.hits }))
    // .then(console.log.bind(console))
  }

  render() {
    return (
      <div className="App">
        <TextField
          id="search"
          name="search"
          value={this.state.search}
          onChange={this.handleChange}
        />
        <Button id="search-button" onClick={this.search}>Search</Button>
        <List>
          {this.state.result.map(item => (
            <ListItem>
              <Card className="result-card">
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
