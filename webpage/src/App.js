import React, { useState, useEffect, useRef } from 'react';
import './App.css';

function App() {
  const [messages, setMessages] = useState([
    /* { content: 'Hello, how can I help you?', author: 'system' },
    { content: 'I need some information about your services.', author: 'user' },
    { content: 'Sure, what would you like to know?', author: 'system' },
    { content: 'Can you tell me about your pricing?', author: 'user' },
    { content: 'Our pricing depends on the services you choose. Please visit our pricing page for more details.', author: 'system' },
    { content: 'Thank you!Thank you!Thank you!Thank you!Thank you!Thank you!Thank you!Thank you!Thank you!Thank you!Thank you!Thank you!', author: 'user' },
    { content: 'You are welcome!', author: 'system' },
    { content: 'Can I get a discount?', author: 'user' },
    { content: 'We offer discounts for bulk orders.', author: 'system' },
    { content: 'How can I place an order?', author: 'user' },
    { content: 'You can place an order through our website.', author: 'system' },
    { content: 'Do you offer customer support?', author: 'user' },
    { content: 'Yes, we offer 24/7 customer support.', author: 'system' },
    { content: 'How can I contact support?', author: 'user' },
    { content: 'You can contact support via email or phone.', author: 'system' },
    { content: 'What is your return policy?', author: 'user' },
    { content: 'We offer a 30-day return policy.', author: 'system' },
    { content: 'Do you ship internationally?', author: 'user' },
    { content: 'Yes, we ship to most countries.', author: 'system' },
    { content: 'How long does shipping take?', author: 'user' },
    { content: 'Shipping times vary depending on the destination.', author: 'system' },
    { content: 'Can I track my order?', author: 'user' },
    { content: 'Yes, you will receive a tracking number once your order is shipped.', author: 'system' },
    { content: 'Thank you for the information.', author: 'user' },
    { content: 'You are welcome!', author: 'system' } */
  ]);
  const [input, setInput] = useState('');
  const messageEndRef = useRef(null);

  const handleSendMessage = () => {
    if (input.trim()) {
      const newMessage = { content: input, author: 'user' };
      const newMessages = [...messages, newMessage];
      setMessages(newMessages);
      setInput('');

      fetch('/api/message', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ messages: newMessages })
      })
        .then(response => response.json())
        .then(data => {
          setMessages([...newMessages, { content: data.data || data.error, author: 'system' }]);
          console.log('Message sent successfully:', data);
        })
        .catch(error => {
          console.error('Error sending message:', error);
        });
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    } else if (e.key === 'Enter' && e.ctrlKey) {
      setInput(input + '\n');
    }
  };

  useEffect(() => {
    messageEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  return (
    <div className="App">
      <header className="App-header">
       Konfig AI Verifier
      </header>
      <div className="message-history">
        {messages.map((message, index) => (
          <div
            key={index}
            className={`message-container ${message.author === 'user' ? 'user-message-container' : 'server-message-container'}`}
          >
            <div className="avatar">
              {message.author === 'user' ? 'You' : 'System'}
            </div>
            <div
              className={`message ${message.author === 'user' ? 'user-message' : 'server-message'}`}
            >{message.content}</div>
          </div>
        ))}
        <div ref={messageEndRef} />
      </div>
      <div className="message-input">
        <textarea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder="Type your message here..."
        />
        <button onClick={handleSendMessage}>Send</button>
      </div>
    </div>
  );
}

export default App;
