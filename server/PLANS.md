🎮 Player & Authentication

POST /auth/register → Create a new account.

POST /auth/login → Authenticate and return session/JWT.

POST /auth/logout → Invalidate session.

GET /player/me → Get current player profile.

PUT /player/me → Update player profile (name, avatar, settings).

GET /player/{id} → Get another player’s public profile.

🏠 Lobby & Matchmaking

POST /lobby/create → Create a lobby (options: map, mode, max players).

POST /lobby/join/{lobbyId} → Join a specific lobby.

POST /lobby/leave/{lobbyId} → Leave a lobby.

GET /lobby/{lobbyId} → Get lobby details and player list.

POST /lobby/start/{lobbyId} → Start match (host or server-controlled).

GET /lobbies → List open lobbies.

POST /matchmaking/join → Join matchmaking queue (ranked/unranked).

POST /matchmaking/cancel → Cancel matchmaking queue.

⚔️ Game Session (Real-Time, via WebSockets)

ws://.../game/connect → Establish real-time connection.

Events you may send/receive over WS:

player.move → Player movement update.

player.action → Attack, interact, use item.

game.state → Broadcast full/partial game state.

game.event → Events like kills, scores, objectives captured.

game.end → Match over, results.

🏆 Leaderboards & Stats

GET /leaderboard/global → Global rankings.

GET /leaderboard/friends → Rankings among friends/guildmates.

GET /leaderboard/{season} → Seasonal rankings.

GET /stats/me → Current player stats.

GET /stats/{playerId} → Another player’s stats.

💰 Inventory & Economy

GET /inventory → List owned items/currency.

POST /inventory/use/{itemId} → Use or equip an item.

POST /inventory/trade → Trade between players.

POST /shop/buy → Buy item (currency/item id in request).

POST /shop/sell → Sell item.

GET /shop → List items available in shop.

💬 Chat & Social

POST /chat/send → Send a chat message (global, lobby, guild).

GET /chat/{channelId} → Fetch recent chat history.

POST /friends/add/{playerId} → Send friend request.

POST /friends/remove/{playerId} → Remove friend.

GET /friends → List friends & their online status.

POST /guild/create → Create a guild/clan.

POST /guild/join/{guildId} → Join a guild.

POST /guild/leave/{guildId} → Leave a guild.

GET /guild/{guildId} → Guild details & member list.

🛡️ Admin & Moderation (optional, but useful)

GET /admin/players → List active players.

POST /admin/kick/{playerId} → Kick a player.

POST /admin/ban/{playerId} → Ban a player.

GET /admin/metrics → Server health (CPU, memory, connections).
