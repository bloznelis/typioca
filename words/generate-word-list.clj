(defn run [[book-uri output-file]]
  (->> (slurp book-uri)
       (re-seq #"[\w|']+")
       (map clojure.string/lower-case)
       (filter #(and (> 10 (count %)) (< 3 (count %))))
       (frequencies)
       (sort-by val)
       (take-last 1000)
       (map key)
       (clojure.string/join "\n")
       (spit output-file)))

(run *command-line-args*)
(println "Done!")
